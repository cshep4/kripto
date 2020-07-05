
default: build

vendor:
	virtualenv ./env
	source env/bin/activate
	pip install boto3
	pip install botocore
	pip install uuid
	pip install pandas
	pip freeze > requirements.txt

build:
	sudo pip3 install -r requirements.txt