
default: build

vendor:
	sudo apt update
	sudo apt install python3-pip
	sudo pip install virtualenv
	virtualenv ./env
	. env/bin/activate
	pip install boto3
	pip install botocore
	pip install uuid
	pip install pandas
	pip freeze > requirements.txt

build:
	sudo apt update
	sudo apt install python3-pip
	sudo pip install virtualenv
	virtualenv ./env
	. env/bin/activate
	pip install boto3
	pip install botocore
	pip install uuid
	pip install pandas
	pip freeze > requirements.txt