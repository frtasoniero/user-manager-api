// MongoDB initialization script for user_management database

// Switch to the user_management database
db = db.getSiblingDB('user_management');

// Create application user with read/write permissions
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

// Create users collection with schema validation
db.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['email', 'password_hash', 'profile', 'created_at', 'updated_at'],
      properties: {
        _id: {
          bsonType: 'string',
          pattern: '^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$',
          description: 'Must be a valid UUID v7'
        },
        email: {
          bsonType: 'string',
          pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',
          description: 'Must be a valid email address'
        },
        password_hash: {
          bsonType: 'string',
          minLength: 1,
          description: 'Must be a non-empty string'
        },
        profile: {
          bsonType: 'object',
          required: ['first_name', 'last_name'],
          properties: {
            first_name: {
              bsonType: 'string',
              minLength: 1
            },
            last_name: {
              bsonType: 'string',
              minLength: 1
            },
            address: {
              bsonType: 'object',
              properties: {
                street: { bsonType: 'string' },
                city: { bsonType: 'string' },
                state: { bsonType: 'string' },
                country: { bsonType: 'string' },
                zip_code: { bsonType: 'string' }
              }
            },
            phone: { bsonType: 'string' },
            birthdate: { bsonType: 'string' },
            nin: { bsonType: 'string' }
          }
        },
        created_at: {
          bsonType: 'date'
        },
        updated_at: {
          bsonType: 'date'
        }
      }
    }
  },
  validationLevel: 'strict',
  validationAction: 'error'
});

// Create indexes for better performance
db.users.createIndex(
  { email: 1 },
  { unique: true, name: 'email_unique_idx' }
);

db.users.createIndex(
  { _id: 1 },
  { name: 'id_time_ordered_idx' }
);

db.users.createIndex(
  { 'profile.first_name': 1, 'profile.last_name': 1 },
  { name: 'name_idx' }
);

db.users.createIndex(
  { 'profile.phone': 1 },
  { sparse: true, name: 'phone_sparse_idx' }
);

db.users.createIndex(
  { 'profile.nin': 1 },
  { unique: true, sparse: true, name: 'nin_unique_sparse_idx' }
);

db.users.createIndex(
  { created_at: 1 },
  { name: 'created_at_idx' }
);

print('✅ Database initialized successfully!');
print('✅ Users collection created with schema validation');
print('✅ Indexes created for optimal performance');
print('✅ Application user created with proper permissions');