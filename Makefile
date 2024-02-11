
build: cron.go
	go build -o suncron cron.go

install: build suncron.conf.example
	cp suncron /usr/local/bin
	test -f /etc/suncron.conf || cp suncron.conf.example /etc/suncron.conf
	test -f /etc/suncron.cron || cp suncron.cron.example /etc/suncron.cron
	touch /etc/cron.d/suncron
	echo "0	0 * * * suncron" > /etc/cron.d/suncron-update
