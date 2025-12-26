package handler

import "net/http"


func CreateJobFromDashboard(w http.ResponseWriter, r *http.Request){
	CreateJob(w,r,true)
}



func CreateJobBatchFromDashboard(w http.ResponseWriter, r *http.Request){
	CreateJob(w,r,true)
}

func CreateJobBatchAtomicFromDashboard(w http.ResponseWriter, r *http.Request){
	CreateJob(w,r,true)
}

func CreateJobFromSystem(w http.ResponseWriter, r *http.Request){
	CreateJob(w,r,false)
}

func CreateJobBatchFromSystem(w http.ResponseWriter, r *http.Request){
	CreateJob(w,r,false)
}

func CreateJobBatchAtomicFromSystem(w http.ResponseWriter, r *http.Request){
	CreateJob(w,r,false)
}