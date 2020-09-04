
default: build

vendor:
	sudo apt update
	sudo apt install python3-pip
	sudo pip install virtualenv
	virtualenv ./env --python=python3
	. env/bin/activate
	python3.7 -m pip install boto3
	python3.7 -m pip install botocore
	python3.7 -m pip install pandas
	python3.7 -m pip freeze > requirements.txt
	cat requirements.txt

build:
	sudo apt update
	sudo apt install python3-pip
	sudo pip install virtualenv
	virtualenv ./env --python=python3
	. env/bin/activate
	python3.7 -m pip install boto3
	python3.7 -m pip install botocore
	python3.7 -m pip install pandas
	python3.7 -m pip freeze > requirements.txt
	cat requirements.txt