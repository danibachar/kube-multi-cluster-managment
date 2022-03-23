
##############################################
#
#
#
#               >>>>>>>
#             /         
#            /
########### /  
# Product # 
########### \ 
#         
#
#
#
#
##############################################
cleanup: delete-bookinfo

build:
	docker build -f Dockerfile.gunicorn -t danibachar/tester:latest .
	docker push danibachar/tester:latest

build-vanila:
	docker build -f Dockerfile -t danibachar/tester-vanila:latest .
	docker push danibachar/tester-vanila:latest

run:
	docker run -p 8080:8081 danibachar/tester:latest

run-vanila:
	docker run -p 8080:8081 danibachar/tester-vanila:latest

deploy-bookinfo:
	kubectl apply -f bookinfo/

delete-bookinfo:
	kubectl delete -f bookinfo/