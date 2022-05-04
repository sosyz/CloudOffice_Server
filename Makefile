all: clear bd
	zip main main scf_bootstrap
clear:
	rm main
bd:
	go build -o main
