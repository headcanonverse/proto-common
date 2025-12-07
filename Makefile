.PHONY: gen
gen:
	protoc -I . --go_out=. --go_opt=module=github.com/headcanonverse/proto-common common/common.proto events/events.proto

.PHONY: clean
clean:
	rm -f common/*.pb.go events/*.pb.go

