#!/bin/bash

docker run -it --network mysql-net --rm mysql mysql -hgo-test-mysql -uroot -pcumdump
