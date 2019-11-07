// Package constant package contains files with various constants in the application.
package constant

// GoogleServiceCredName supplies the name of the file containing credentials for our
// Google service account. Used to authenticate with firebase.
const GoogleServiceCredName = "gt-backend-12f5a-firebase-adminsdk-agjnd-2b80155741.json"

// GoogleCredEnvVar - name of environment variable for hosting cred file in cloud
const GoogleCredEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

// FirebaseStorageBucket url
const FirebaseStorageBucket = "gt-backend-12f5a.appspot.com"

// IgcMetadata - Name of the collection containing igc metadata in firestore
const IgcMetadata = "igc-metadata"

// PageSize - Size of a page in a Firestore Query
const PageSize = 20

// FirebaseQueryOrder - What to order the Firestore Query with (time)
const FirebaseQueryOrder = "Time"

// MaxIgcFileSize - Max size of an IGC file (MB - KB - B)
const MaxIgcFileSize = 15 * 1024 * 1024

// ScraperUID - Scraper UID
const ScraperUID = "iP1dgAHJ2JNce4hGr9H0RugkCHP2"

// TestUID - E2E Test UID
const TestUID = "o1Sz791YSHby0PCe51JlxSD6G533"
