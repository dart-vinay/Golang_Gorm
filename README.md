# Golang API (GORM, goFiber, Postgres)

## About

The repository provides basic introduction of gorm and goFiber(web framework of golang) with postgres as supporting database.
This repository has four main folders and a few data generation scripts:

* **data**: Contains all the raw data we need to populate our elastic search indexes
* **routes**: Contains the API definition and corresponding request routes
* **models**: Contains the structural objects used in this project and implements all the logic used
* **util**: Contains all the helper function that we need for functioning of our services

## Installation

Install Golang, GORM, GoFiber, PostgreSQL

## Usage


```bash
go run main.go
```
This will start the backend server at `localhost:3000`

To test the behavior out we have four APIs:
* Fetch all the sheets for user with marks and comments corresponding to each question
  * endPoint: `/v1/sheet/student`
  * Request Type : **GET**
  * Request Headers:
    * X-USER-ID : < USER ID > 
    * X-USER-TYPE : <USER TYPE (STUDENT, TEACHER)>
  * Response: 200 Ok
 ```javascript
{
    "sheet_data": [
        {
            "SheetId": 1,
            "Questions": [
                {
                    "QuestionId": 1,
                    "QuestionRefId": 5,
                    "QuestionText": "Question with ID 1",
                    "CorrectedBy": 2,
                    "MaxMarks": 10,
                    "MarksAlloted": 7,
                    "Comments": [
                        {
                            "ID": 3,
                            "CommentText": "User 2 Commented",
                            "CommentStatus": "PENDING",
                            "UserId": 2,
                            "QuestionUserReferenceId": 5,
                            "CreatedAt": "2021-08-26T16:23:03.479968+05:30",
                            "UpdatedAt": "2021-08-26T16:23:03.479968+05:30",
                            "DeletedAt": null
                        }
                    ]
                },
                {
                    "QuestionId": 2,
                    "QuestionRefId": 6,
                    "QuestionText": "Question with ID 2",
                    "CorrectedBy": 2,
                    "MaxMarks": 20,
                    "MarksAlloted": 10,
                    "Comments": null
                }
            ],
            "Subject": "ENGLISH"
        }
    ]
}
```
* Post comments corresponding to question which needs to be remarked
  * endPoint: `/v1/sheet/student/comment`
  * Request Type : **POST**
   * Request Headers:
      * X-USER-ID : < USER ID >
      * X-USER-TYPE : <USER TYPE (STUDENT, TEACHER)>
    * Response: 200 Ok
  * Request Object :
```javascript
{
    "user_comments" : [
        {
            "question_ref_id" : 3,
            "comment" : "User 1 commented"
        },
        {
            "question_ref_id" : 5,
            "comment" : "User 2 Commented"
        }
    ]
}
```  
* Fetch all the relevant comments for teacher
  * endPoint: `/v1/teacher/comment`
  * Request Type : **GET**
  * Request Headers:
    * X-USER-ID : < USER ID > 
    * X-USER-TYPE : <USER TYPE (STUDENT, TEACHER)>
  * Response: 200 Ok
 ```javascript
{
    "comments": [
        {
            "ID": 4,
            "CommentText": "User 1 commented",
            "CommentStatus": "REMARKED",
            "UserId": 1,
            "QuestionUserReferenceId": 3,
            "CreatedAt": "2021-08-26T16:24:29.279059+05:30",
            "UpdatedAt": "2021-08-26T16:25:52.672309+05:30",
            "DeletedAt": null
        }
    ]
}
```
* Action taken by teacher for their comments
  * endPoint: `/v1/teacher/comment/action`
  * Request Type : **POST**
   * Request Headers:
      * X-USER-ID : < USER ID >
      * X-USER-TYPE : <USER TYPE (STUDENT, TEACHER)>
    * Response: 200 Ok
  * Request Object :
```javascript
{
    "actions" : [
        {
           "question_ref_id" : 3,
            "comment_id" : 4,
            "comment_action" : "REMARKED",
            "modified_marks" : 5  
        }
    ]
}
```  
Go ahead and test it and let me know your thoughts over this.

## Updates

Will be grooming this repository more with further developments. Cheers!