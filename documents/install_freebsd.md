# Install FreeBSD

```bash
pkg install mongodb80
pkg install node
pkg install npm
pkg install mongodb-tools
```

mongodb 6.0 이상부터는 mongo 명령어가 존재하지 않는다. 따로 설치한다.

```bash
$ npm install mongosh
$ npx mongosh

>> show dbs
```

## Install Font

```bash
sudo pkg install -y freefont-ttf
```

FreeMono Path

```
/usr/local/share/fonts/freefont-ttf/FreeMono.ttf
```

## Install mongodb50-5.0.29

```
sudo pkg install mongodb50
sudo sysrc mongod_enable="YES"
mongo
```
