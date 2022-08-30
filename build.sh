rm -rf .product
mkdir .product
mkdir .product/static

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go build -o go-fly-pro

cp go-fly-pro .product/
cp -r config  .product/config
cp -r static/js  .product/static/js
cp -r static/css  .product/static/css
cp -r static/images  .product/static/images
cp -r static/templates .product/static/templates
cp help.txt .product/
cp stop.sh .product/
cd .product/ && zip -q -r go-fly-pro.zip   *