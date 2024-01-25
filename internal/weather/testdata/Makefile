include .env

all:				\
	creston,ca.json		\
	victoria,ca.json	\
	vancouver,ca.json	\
	toronto,ca.json		\
	chicago,us.json		\
	gibsons,ca.json		\
	shanghai,cn.json

%.json:
	wget "http://api.openweathermap.org/data/2.5/weather?units=metric&q=$*&appid=${OPENWEATHERMAP_API_KEY}" -O $@

