db.users.insertMany([
  {
    "_id": ObjectId("5d470b3e98b0116d7d8ca48c"),
    "email": "test1@test.com",
    "password": "$2a$10$Fl2qgZ7DjYarrymLT6tLle3CqQ.LLdQ/U1E2XCvB6tqFwN4Q5m09a",
    "roles": ["user"],
    "active": true,
    "created_at": ISODate("2019-08-07T00:58:07.162Z"),
    "updated_at": ISODate("2019-08-07T00:58:07.162Z")
  },
  {
    "_id": ObjectId("5d4a22b1106eded67d47c02e"),
    "email": "test2@test.com",
    "password": "$2a$10$Fl2qgZ7DjYarrymLT6tLle3CqQ.LLdQ/U1E2XCvB6tqFwN4Q5m09a",
    "roles": ["user", "sysadmin"],
    "active": false,
    "created_at": ISODate("2019-08-07T00:58:07.162Z"),
    "updated_at": ISODate("2019-08-07T00:58:07.162Z")
  },
  {
    "_id": ObjectId("5d4a22e9587f3dbb8d33fd38"),
    "email": "test3@test.com",
    "password": "$2a$10$Fl2qgZ7DjYarrymLT6tLle3CqQ.LLdQ/U1E2XCvB6tqFwN4Q5m09a",
    "roles": ["user"],
    "active": true,
    "created_at": ISODate("2019-08-07T00:58:07.162Z"),
    "updated_at": ISODate("2019-08-07T00:58:07.162Z")
  },
  {
    "_id": ObjectId("5d4a22e9587f3dbb8d33fd39"),
    "email": "test4@test.com",
    "password": "$2a$10$Fl2qgZ7DjYarrymLT6tLle3CqQ.LLdQ/U1E2XCvB6tqFwN4Q5m09a",
    "roles": ["user"],
    "active": true,
    "created_at": ISODate("2019-08-07T00:58:07.162Z"),
    "updated_at": ISODate("2019-08-07T00:58:07.162Z")
  }
]);