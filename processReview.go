package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/alfg/mp4"
	"github.com/amarburg/go-quicktime"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProcessReviewRender 함수는 OpenPipelineIO 가 실행되면서 처리될 프로세싱을 진행한다.
func ProcessReviewRender() {
	// 버퍼 채널을 만든다.
	jobs := make(chan Review, *flagProcessBufferSize) // 작업을 대기할 버퍼를 만든다.
	// worker 프로세스를 지정한 개수만큼 실행시킨다.
	for w := 1; w <= *flagMaxProcessNum; w++ {
		go workerReview(jobs)
	}
	// queueingItem을 실행시킨다.
	go queueingReviewItem(jobs)

	// ProcessMain()이 종료되지 않는 역할을 한다.
	select {}
}

// workerReview 합수는 Review 데이터를 jobs로 보낸다.
func workerReview(jobs <-chan Review) {
	for job := range jobs {
		// job은 리뷰타입이다.
		switch job.Type {
		case "image":
			processingReviewImageItem(job)
		default:
			processingReviewClipItem(job)
		}
	}
}

// queueingReviewItem 은 연산할 Review 아이템을 jobs 채널에 전송한다.
func queueingReviewItem(jobs chan<- Review) {
	for {
		if *flagDebug {
			fmt.Printf("wait %d sec before review process\n", *flagProcessDuration)
		}
		time.Sleep(time.Second * time.Duration(*flagProcessDuration))
		// ProcessStatus가 wait인 item을 가져온다.
		review, err := GetWaitProcessStatusReview() // 이 함수로 반환되는 아이템은 리뷰 아이템은 상태가 queued가 된 리뷰 아이템이다.
		if err != nil {
			// 가지고 올 문서가 없다면 기다렸다가 continue.
			if err == mongo.ErrNoDocuments {
				continue
			}
			continue
		}
		if *flagDebug {
			fmt.Println(review)
		}
		jobs <- review
	}
}

func processingReviewClipItem(review Review) {
	client, err := initMongoClient()
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Disconnect(context.Background())
	reviewID := review.ID.Hex()
	// 연산 상태를 queued 으로 바꾼다. 바꾸는 이유는 ffmpeg 연산이 10초이상 진행될 때 상태가 바뀌지 않아서 이전에 연산중인 데이터가 다시 연산될 수 있기 때문이다.
	err = setReviewProcessStatusV2(client, reviewID, "processing")
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	// ffmpeg 경로를 체크한다.
	if _, err := os.Stat(CachedAdminSetting.FFmpeg); os.IsNotExist(err) {
		err = setErrReviewV2(client, reviewID, "ffmpeg가 존재하지 않습니다")
		if err != nil {
			log.Println(err)
		}
		return
	}
	// ReviewDataPath가 존재하는지 경로를 체크한다.
	if _, err := os.Stat(CachedAdminSetting.ReviewDataPath); os.IsNotExist(err) {
		err = setErrReviewV2(client, reviewID, "admin 셋팅에 ReviewDataPath가 존재하지 않습니다")
		if err != nil {
			log.Println(err)
		}
		return
	}
	// review데이터가 atom 구조를 같는지 체크한다.
	err = checkQuicktimeFileStruct(review)
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	// mp4를 생성한다.
	err = genMp4(review)
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	// 생성된 .mp4 파일이 mp4 자료구조를 같는지 체크한다.
	err = checkMp4FileStruct(review)
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	// 연산이 끝나고 해당 파일을 삭제해야 한다면 삭제를 진행한다.
	if review.RemoveAfterProcess {
		err = os.Remove(review.Path)
		if err != nil {
			err = setErrReviewV2(client, reviewID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	// 연산 상태를 done 으로 바꾼다.
	err = setReviewProcessStatusV2(client, reviewID, "done")

	// 리뷰 데이터를 추가하고 나서 "앗, 리뷰 잘못올렸네~ 취소해야지~하면서.." 서버에서 연산중인 리뷰데이터를 바로 삭제하는 아티스트가 있다.
	// 이러한 상황에서는 삭제가 일어나면 상태를 바꿀 review 아이템이 DB에 없게 된다.
	// 만약 상태를 바꾸어야 할 때 해당 리뷰아이템이 없다면 바로 return 하도록 하였다.
	if err == mongo.ErrNoDocuments {
		return
	}
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func processingReviewImageItem(review Review) {
	client, err := initMongoClient()
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Disconnect(context.Background())
	reviewID := review.ID.Hex()
	// 연산 상태를 queued 으로 바꾼다. 바꾸는 이유는 ffmpeg 연산이 10초이상 진행될 때 상태가 바뀌지 않아서 이전에 연산중인 데이터가 다시 연산될 수 있기 때문이다.
	err = setReviewProcessStatusV2(client, reviewID, "processing")
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	// ReviewDataPath가 존재하는지 경로를 체크한다.
	if _, err := os.Stat(CachedAdminSetting.ReviewDataPath); os.IsNotExist(err) {
		err = setErrReviewV2(client, reviewID, "admin 셋팅에 ReviewDataPath가 존재하지 않습니다")
		if err != nil {
			log.Println(err)
		}
		return
	}
	// image를 리뷰폴더에 복사한다.
	input, err := os.ReadFile(review.Path)
	if err != nil {
		log.Println(err)
		return
	}
	per, err := strconv.ParseInt(CachedAdminSetting.ReviewDataPathPermission, 8, 64)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.WriteFile(CachedAdminSetting.ReviewDataPath+"/"+reviewID+review.Ext, input, os.FileMode(per))
	if err != nil {
		log.Println(err)
		return
	}

	// 이미지 연산된 경로가 review 자료구조의 Path값에 들어가야 한다. 업로드된 이미지 경로는 삭제될 수 있기 때문이다.
	err = setReviewPathV2(client, reviewID, CachedAdminSetting.ReviewDataPath+"/"+reviewID+review.Ext)
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 연산이 끝나고 해당 파일을 삭제해야 한다면 삭제를 진행한다.
	if review.RemoveAfterProcess {
		err = os.Remove(review.Path)
		if err != nil {
			err = setErrReviewV2(client, reviewID, err.Error())
			if err != nil {
				log.Println(err)
			}
			return
		}
	}

	// 연산 상태를 done 으로 바꾼다.
	err = setReviewProcessStatusV2(client, reviewID, "done")
	if err != nil {
		err = setErrReviewV2(client, reviewID, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
}

// checkQuicktimeFileStruct 함수는 리뷰 아이템 정보를 이용해서 atom 구조가 정상인지 체크한다.
func checkQuicktimeFileStruct(item Review) error {
	file, err := os.Open(item.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	_, err = quicktime.BuildTree(file, uint64(info.Size()))
	if err != nil {
		return err
	}
	return nil
}

// checkMp4FileStruct 함수는 리뷰 아이템 정보를 mp4 구조가 정상인지 체크한다.
func checkMp4FileStruct(item Review) error {
	file, err := os.Open(CachedAdminSetting.ReviewDataPath + "/" + item.ID.Hex() + ".mp4")
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	_, err = mp4.OpenFromReader(file, info.Size())
	if err != nil {
		return err
	}
	return nil
}

// genMp4 는 리뷰 아이템 정보를 이용해서 .mp4 동영상을 만든다.
func genMp4(item Review) error {
	args := []string{
		"-y",
		"-i",
		item.Path,
		"-c:v",
		"libx264",
		"-qscale:v",
		"7",
		"-vf",
		"pad=ceil(iw/2)*2:ceil(ih/2)*2", // 영상의 세로 픽셀이 홀수일때 연산되지 않는다. 이 옵션이 필요하다.
		"-pix_fmt",
		"yuv420p", // 이 옵션이 없다면 Prores로 동영상을 만들때 크롬에서만 재생된다.

	}
	if CachedAdminSetting.AudioCodec == "nosound" {
		// nosound라면 사운드를 넣지 않는 옵션을 추가한다.
		args = append(args, "-an")
	} else {
		// 다른 사운드 코덱이라면 사운드클 체크한다.
		args = append(args, "-c:a")
		args = append(args, CachedAdminSetting.AudioCodec)
	}

	args = append(args, "-threads")
	args = append(args, strconv.Itoa(CachedAdminSetting.FFmpegThreads)) // 웹서버의 부하를 줄이기 위해서 서버수가 적다면 쓰레드 1개만 사용한다.
	args = append(args, CachedAdminSetting.ReviewDataPath+"/"+item.ID.Hex()+".mp4")
	err := exec.Command(CachedAdminSetting.FFmpeg, args...).Run()
	if err != nil {
		return err
	}
	return nil
}
