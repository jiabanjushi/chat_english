rd .product /s/q
mkdir .product
go build -o go-fly-pro.exe
xcopy  config  .product\config /e /i /s /y
xcopy  static\images  .product\static\images /e /i /s /y
xcopy  static\css  .product\static\css /e /i /s /y
xcopy  static\js  .product\static\js /e /i /s /y
copy go-fly-pro.exe .product\  /y
copy cron.bat .product\  /y
copy start.bat .product\  /y
copy stop.bat .product\  /y
copy install.bat .product\  /y
copy help.txt .product\  /y
del go-fly-pro.exe
pause