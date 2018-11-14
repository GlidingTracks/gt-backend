// Package constant package contains files with various constants in the application.
package constant

// GoogleServiceCredName supplies the name of the file containing credentials for our
// Google service account. Used to authenticate with firebase.
const GoogleServiceCredName = "gt-backend-8b9c2-firebase-adminsdk-0t965-d5b53ac637.json"

// GoogleCredEnvVar - name of environment variable for hosting cred file in cloud
const GoogleCredEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

// FirebaseStorageBucket url
const FirebaseStorageBucket = "gt-backend-8b9c2.appspot.com"

// IgcMetadata - Name of the collection containing igc metadata in firestore
const IgcMetadata = "igc-metadata"

// PageSize - Size of a page in a Firestore Query
const PageSize = 20

// FirebaseQueryOrder - What to order the Firestore Query with (time)
const FirebaseQueryOrder = "Time"

// MaxIgcFileSize - Max size of an IGC file (MB - KB - B)
const MaxIgcFileSize = 15 * 1024 * 1024

// ScraperUID - Scraper UID
const ScraperUID = "Hvl8FlvVz6QSKkNjsDeLB5aKBmA2"

// TestUID - E2E Test UID
const TestUID = "o1Sz791YSHby0PCe51JlxSD6G533"
