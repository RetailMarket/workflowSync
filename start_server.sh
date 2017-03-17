if [[ ! -e out/ ]];
		then
			mkdir out/
		fi 

		go build -o out/build app/workflowSync/main.go; ./out/build
