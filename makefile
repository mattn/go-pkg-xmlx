
all:
	make -C xmlx install

test:
	make -C xmlx test

clean:
	make -C xmlx clean

format:
	gofmt -w .
