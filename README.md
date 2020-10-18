# adminer_api
Simple CLI api for AdminerPHP

You can export databases from Adminer.
```
go run main.go -url https://example.com/adminer.php -pass admin123 -user admin -exportOutput save                     
```

For more info:
```
go run main.php -help
```

To do:

1) Select db for export
2) Add proxy support
3) SQL shell
