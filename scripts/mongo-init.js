// MongoDB initialization script
// This script runs when the container starts for the first time

// Switch to the user_management database
db = db.getSiblingDB('user_management');

// Create a user for the application
db.createUser({
  user: 'api_user',
  pwd: 'api_password',
  roles: [
    {
      role: 'readWrite',
      db: 'user_management'
    }
  ]
});

// Create users collection with validation
db.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['email', 'password_hash', 'profile', 'created_at', 'updated_at'],
      properties: {
        _id: {
          bsonType: 'objectId',
          description: 'MongoDB ObjectID'
        },
        email: {
          bsonType: 'string',
          pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
          description: 'Must be a valid email address'
        },
        password_hash: {
          bsonType: 'string',
          minLength: 1,
          description: 'Must be a hashed password'
        },
        profile: {
          bsonType: 'object',
          properties: {
            first_name: {
              bsonType: 'string',
              description: 'User first name'
            },
            last_name: {
              bsonType: 'string',
              description: 'User last name'
            },
            address: {
              bsonType: 'object',
              properties: {
                street: {
                  bsonType: 'string',
                  description: 'Street address'
                },
                city: {
                  bsonType: 'string',
                  description: 'City name'
                },
                state: {
                  bsonType: 'string',
                  description: 'State or province'
                },
                country: {
                  bsonType: 'string',
                  description: 'Country name'
                },
                zip_code: {
                  bsonType: 'string',
                  description: 'Postal/ZIP code'
                }
              },
              additionalProperties: false
            },
            phone: {
              bsonType: 'string',
              description: 'Phone number'
            },
            birthdate: {
              bsonType: 'string',
              description: 'Birth date as string (YYYY-MM-DD format recommended)'
            },
            nin: {
              bsonType: 'string',
              description: 'National identification number'
            }
          },
          additionalProperties: false
        },
        created_at: {
          bsonType: 'date',
          description: 'Document creation timestamp'
        },
        updated_at: {
          bsonType: 'date',
          description: 'Document last update timestamp'
        }
      },
      additionalProperties: false
    }
  }
});

// Create unique index on email
db.users.createIndex({ email: 1 }, { unique: true });

// Create index on created_at for efficient sorting
db.users.createIndex({ created_at: 1 });

// Create index on profile fields for searching
db.users.createIndex({ 
  "profile.first_name": 1, 
  "profile.last_name": 1 
});

// Create index on phone for lookups (sparse index since it's optional)
db.users.createIndex({ "profile.phone": 1 }, { sparse: true });

// Create index on NIN for lookups (sparse and unique since it's optional but should be unique)
db.users.createIndex({ "profile.nin": 1 }, { sparse: true, unique: true });

print('Database initialization completed successfully!');
