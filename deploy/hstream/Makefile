
sql:
	docker run -it --rm --name some-hstream-cli-1 -v /Users/gangtao/Code/personal/chameleon/deploy/hstream/dbdata:/data/store --network host hstreamdb/hstream hstream-client --port 6570

store:
	docker run -td --rm --name some-hstream-store -v /Users/gangtao/Code/personal/chameleon/deploy/hstream/dbdata:/data/store --network host hstreamdb/hstream ld-dev-cluster --root /data/store --use-tcp

server:
	docker run -it --rm --name some-hstream-server -v /Users/gangtao/Code/personal/chameleon/deploy/hstream/dbdata:/data/store --network host hstreamdb/hstream hstream-server --port 6570 --store-config /data/store/logdevice.conf