#!/bin/bash

go build -o bookings cmd/web/*.go
./bookings  #passing flags here  (database info and others)
