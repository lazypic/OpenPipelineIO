# Install FreeBSD

```bash
pkg install mongodb70
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

## Install mongodb50-5.0.29

```
sudo pkg install mongodb50
mongo
```
