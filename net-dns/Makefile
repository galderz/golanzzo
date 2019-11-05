build:
	go build -o nslookup-go

clean:
	rm nslookup-go

image:
	docker build -t galderz/nslookup-go .

push:
	docker push galderz/nslookup-go

run:
	oc run nslookup-go --image=galderz/nslookup-go --restart=Never --attach -i --tty
