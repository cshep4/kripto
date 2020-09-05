#
#default: build
#
#vendor:
#	sudo apt update
#	sudo apt install python3-pip
#	sudo pip install virtualenv
#	virtualenv ./env --python=python3
#	. env/bin/activate
#	python3.7 -m pip install boto3
#	python3.7 -m pip install botocore
#	python3.7 -m pip install pandas
#	python3.7 -m pip freeze > requirements.txt
#	cat requirements.txt
#
#build:
#	sudo apt update
#	sudo apt install python3-pip
#	sudo pip install virtualenv
#	virtualenv ./env --python=python3
#	. env/bin/activate
#	python3.7 -m pip install boto3
#	python3.7 -m pip install botocore
#	python3.7 -m pip install pandas
#	python3.7 -m pip freeze > requirements.txt
#	cat requirements.txt

default: build

vendor:
	yum -y install sudo
	yum -y install synaptic
	sudo apt update
	sudo apt install python3-pip
	sudo pip install virtualenv
	virtualenv ./env
	. env/bin/activate
	pip install --target ./package boto3
	pip install --target ./package botocore
	pip install --target ./package numpy
	pip install --target ./package pandas
	pip freeze > requirements.txt

build:
#	virtualenv ./env
#	. env/bin/activate
	python3.7 -m pip install -t . boto3
	#python3.7 -m pip install -t . botocore
	python3.7 -m pip install -t . pytz
	rm -r *.dist-info __pycache__
	#rm -r *.whl *.dist-info __pycache__
	#pip freeze > requirements.txt

create:
	docker build -t lambda_image .
	docker create --name extract lambda_image
	docker cp extract:/var/task/lambda.zip .
	docker rm extract