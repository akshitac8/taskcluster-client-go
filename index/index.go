// The following code is AUTO-GENERATED. Please DO NOT edit.
// To update this generated code, run the following command:
// in the /codegenerator/model subdirectory of this project,
// making sure that `${GOPATH}/bin` is in your `PATH`:
//
// go install && go generate
//
// This package was generated from the schema defined at
// http://references.taskcluster.net/index/v1/api.json

// The task index, typically available at `index.taskcluster.net`, is
// responsible for indexing tasks. The service ensures that tasks can be
// located by recency and/or arbitrary strings. Common use-cases include:
//
//  * Locate tasks by git or mercurial `<revision>`, or
//  * Locate latest task from given `<branch>`, such as a release.
//
// **Index hierarchy**, tasks are indexed in a dot (`.`) separated hierarchy
// called a namespace. For example a task could be indexed with the index path
// `some-app.<revision>.linux-64.release-build`. In this case the following
// namespaces is created.
//
//  1. `some-app`,
//  1. `some-app.<revision>`, and,
//  2. `some-app.<revision>.linux-64`
//
// Inside the namespace `some-app.<revision>` you can find the namespace
// `some-app.<revision>.linux-64` inside which you can find the indexed task
// `some-app.<revision>.linux-64.release-build`. This is an example of indexing
// builds for a given platform and revision.
//
// **Task Rank**, when a task is indexed, it is assigned a `rank` (defaults
// to `0`). If another task is already indexed in the same namespace with
// lower or equal `rank`, the index for that task will be overwritten. For example
// consider index path `mozilla-central.linux-64.release-build`. In
// this case one might choose to use a UNIX timestamp or mercurial revision
// number as `rank`. This way the latest completed linux 64 bit release
// build is always available at `mozilla-central.linux-64.release-build`.
//
// Note that this does mean index paths are not immutable: the same path may
// point to a different task now than it did a moment ago.
//
// **Indexed Data**, when a task is retrieved from the index the result includes
// a `taskId` and an additional user-defined JSON blob that was indexed with
// the task.
//
// **Entry Expiration**, all indexed entries must have an expiration date.
// Typically this defaults to one year, if not specified. If you are
// indexing tasks to make it easy to find artifacts, consider using the
// artifact's expiration date.
//
// **Valid Characters**, all keys in a namespace `<key1>.<key2>` must be
// in the form `/[a-zA-Z0-9_!~*'()%-]+/`. Observe that this is URL-safe and
// that if you strictly want to put another character you can URL encode it.
//
// **Indexing Routes**, tasks can be indexed using the API below, but the
// most common way to index tasks is adding a custom route to `task.routes` of the
// form `index.<namespace>`. In order to add this route to a task you'll
// need the scope `queue:route:index.<namespace>`. When a task has
// this route, it will be indexed when the task is **completed successfully**.
// The task will be indexed with `rank`, `data` and `expires` as specified
// in `task.extra.index`. See the example below:
//
// ```js
// {
//   payload:  { /* ... */ },
//   routes: [
//     // index.<namespace> prefixed routes, tasks CC'ed such a route will
//     // be indexed under the given <namespace>
//     "index.mozilla-central.linux-64.release-build",
//     "index.<revision>.linux-64.release-build"
//   ],
//   extra: {
//     // Optional details for indexing service
//     index: {
//       // Ordering, this taskId will overwrite any thing that has
//       // rank <= 4000 (defaults to zero)
//       rank:       4000,
//
//       // Specify when the entries expire (Defaults to 1 year)
//       expires:          new Date().toJSON(),
//
//       // A little informal data to store along with taskId
//       // (less 16 kb when encoded as JSON)
//       data: {
//         hgRevision:   "...",
//         commitMessae: "...",
//         whatever...
//       }
//     },
//     // Extra properties for other services...
//   }
//   // Other task properties...
// }
// ```
//
// **Remark**, when indexing tasks using custom routes, it's also possible
// to listen for messages about these tasks. For
// example one could bind to `route.index.some-app.*.release-build`,
// and pick up all messages about release builds. Hence, it is a
// good idea to document task index hierarchies, as these make up extension
// points in their own.
//
// See: https://docs.taskcluster.net/reference/core/index/api-docs
//
// How to use this package
//
// First create an Index object:
//
//  myIndex, err := index.New(nil)
//
// handling any errors...
//
//  if err != nil {
//  	// handle error...
//  }
//
// and then call one or more of myIndex's methods, e.g.:
//
//  data, err := myIndex.FindTask(.....)
//
// handling any errors...
//
//  if err != nil {
//  	// handle error...
//  }
//
// Taskcluster Schema
//
// The source code of this go package was auto-generated from the API definition at
// http://references.taskcluster.net/index/v1/api.json together with the input and output schemas it references, downloaded on
// Mon, 19 Mar 2018 at 14:21:00 UTC. The code was generated
// by https://github.com/taskcluster/taskcluster-client-go/blob/master/build.sh.
package index

import (
	"net/url"
	"time"

	tcclient "github.com/taskcluster/taskcluster-client-go"
)

const (
	DefaultBaseURL = "https://index.taskcluster.net/v1"
)

type Index tcclient.Client

// New returns an Index client, configured to run against production. Pass in
// nil to load credentials from TASKCLUSTER_* environment variables. The
// returned client is mutable, so returned settings can be altered.
//
//  myIndex, err := index.New(nil)                           // credentials loaded from TASKCLUSTER_* env vars
//  if err != nil {
//      // handle malformed credentials...
//  }
//  myIndex.BaseURL = "http://localhost:1234/api/Index/v1"   // alternative API endpoint (production by default)
//  data, err := myIndex.FindTask(.....)                     // for example, call the FindTask(.....) API endpoint (described further down)...
//  if err != nil {
//  	// handle errors...
//  }
//
// If authentication is not required, use NewNoAuth() instead.
func New(credentials *tcclient.Credentials) (*Index, error) {
	if credentials == nil {
		credentials = tcclient.CredentialsFromEnvVars()
	}
	err := credentials.Validate()
	myIndex := Index(tcclient.Client{
		Credentials:  credentials,
		BaseURL:      DefaultBaseURL,
		Authenticate: true,
	})
	return &myIndex, err
}

// NewNoAuth returns an Index client with authentication disabled. This is
// useful when calling taskcluster APIs that do not require authorization.
func NewNoAuth() *Index {
	myIndex := Index(tcclient.Client{
		BaseURL:      DefaultBaseURL,
		Authenticate: false,
	})
	return &myIndex
}

// Find a task by index path, returning the highest-rank task with that path. If no
// task exists for the given path, this API end-point will respond with a 404 status.
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#findTask
func (myIndex *Index) FindTask(indexPath string) (*IndexedTaskResponse, error) {
	cd := tcclient.Client(*myIndex)
	responseObject, _, err := (&cd).APICall(nil, "GET", "/task/"+url.QueryEscape(indexPath), new(IndexedTaskResponse), nil)
	return responseObject.(*IndexedTaskResponse), err
}

// List the namespaces immediately under a given namespace.
//
// This endpoint
// lists up to 1000 namespaces. If more namespaces are present, a
// `continuationToken` will be returned, which can be given in the next
// request. For the initial request, the payload should be an empty JSON
// object.
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#listNamespaces
func (myIndex *Index) ListNamespaces(namespace, continuationToken, limit string) (*ListNamespacesResponse, error) {
	v := url.Values{}
	if continuationToken != "" {
		v.Add("continuationToken", continuationToken)
	}
	if limit != "" {
		v.Add("limit", limit)
	}
	cd := tcclient.Client(*myIndex)
	responseObject, _, err := (&cd).APICall(nil, "GET", "/namespaces/"+url.QueryEscape(namespace), new(ListNamespacesResponse), v)
	return responseObject.(*ListNamespacesResponse), err
}

// Stability: *** DEPRECATED ***
//
// List the namespaces immediately under a given namespace.
//
// This endpoint
// lists up to 1000 namespaces. If more namespaces are present, a
// `continuationToken` will be returned, which can be given in the next
// request. For the initial request, the payload should be an empty JSON
// object.
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#listNamespacesPost
func (myIndex *Index) ListNamespacesPost(namespace string, payload *ListNamespacesRequest) (*ListNamespacesResponse, error) {
	cd := tcclient.Client(*myIndex)
	responseObject, _, err := (&cd).APICall(payload, "POST", "/namespaces/"+url.QueryEscape(namespace), new(ListNamespacesResponse), nil)
	return responseObject.(*ListNamespacesResponse), err
}

// List the tasks immediately under a given namespace.
//
// This endpoint
// lists up to 1000 tasks. If more tasks are present, a
// `continuationToken` will be returned, which can be given in the next
// request. For the initial request, the payload should be an empty JSON
// object.
//
// **Remark**, this end-point is designed for humans browsing for tasks, not
// services, as that makes little sense.
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#listTasks
func (myIndex *Index) ListTasks(namespace, continuationToken, limit string) (*ListTasksResponse, error) {
	v := url.Values{}
	if continuationToken != "" {
		v.Add("continuationToken", continuationToken)
	}
	if limit != "" {
		v.Add("limit", limit)
	}
	cd := tcclient.Client(*myIndex)
	responseObject, _, err := (&cd).APICall(nil, "POST", "/tasks/"+url.QueryEscape(namespace), new(ListTasksResponse), v)
	return responseObject.(*ListTasksResponse), err
}

// Stability: *** DEPRECATED ***
//
// List the tasks immediately under a given namespace.
//
// This endpoint
// lists up to 1000 tasks. If more tasks are present, a
// `continuationToken` will be returned, which can be given in the next
// request. For the initial request, the payload should be an empty JSON
// object.
//
// **Remark**, this end-point is designed for humans browsing for tasks, not
// services, as that makes little sense.
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#listTasksPost
func (myIndex *Index) ListTasksPost(namespace string, payload *ListTasksRequest) (*ListTasksResponse, error) {
	cd := tcclient.Client(*myIndex)
	responseObject, _, err := (&cd).APICall(payload, "POST", "/tasks/"+url.QueryEscape(namespace), new(ListTasksResponse), nil)
	return responseObject.(*ListTasksResponse), err
}

// Insert a task into the index.  If the new rank is less than the existing rank
// at the given index path, the task is not indexed but the response is still 200 OK.
//
// Please see the introduction above for information
// about indexing successfully completed tasks automatically using custom routes.
//
// Required scopes:
//   index:insert-task:<namespace>
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#insertTask
func (myIndex *Index) InsertTask(namespace string, payload *InsertTaskRequest) (*IndexedTaskResponse, error) {
	cd := tcclient.Client(*myIndex)
	responseObject, _, err := (&cd).APICall(payload, "PUT", "/task/"+url.QueryEscape(namespace), new(IndexedTaskResponse), nil)
	return responseObject.(*IndexedTaskResponse), err
}

// Find a task by index path and redirect to the artifact on the most recent
// run with the given `name`.
//
// Note that multiple calls to this endpoint may return artifacts from differen tasks
// if a new task is inserted into the index between calls. Avoid using this method as
// a stable link to multiple, connected files if the index path does not contain a
// unique identifier.  For example, the following two links may return unrelated files:
// * https://index.taskcluster.net/task/some-app.win64.latest.installer/artifacts/public/installer.exe`
// * https://index.taskcluster.net/task/some-app.win64.latest.installer/artifacts/public/debug-symbols.zip`
//
// This problem be remedied by including the revision in the index path or by bundling both
// installer and debug symbols into a single artifact.
//
// If no task exists for the given index path, this API end-point responds with 404.
//
// Required scopes:
//   If private:
//     queue:get-artifact:<name>
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#findArtifactFromTask
func (myIndex *Index) FindArtifactFromTask(indexPath, name string) error {
	cd := tcclient.Client(*myIndex)
	_, _, err := (&cd).APICall(nil, "GET", "/task/"+url.QueryEscape(indexPath)+"/artifacts/"+url.QueryEscape(name), nil, nil)
	return err
}

// Returns a signed URL for FindArtifactFromTask, valid for the specified duration.
//
// Required scopes:
//   If private:
//     queue:get-artifact:<name>
//
// See FindArtifactFromTask for more details.
func (myIndex *Index) FindArtifactFromTask_SignedURL(indexPath, name string, duration time.Duration) (*url.URL, error) {
	cd := tcclient.Client(*myIndex)
	return (&cd).SignedURL("/task/"+url.QueryEscape(indexPath)+"/artifacts/"+url.QueryEscape(name), nil, duration)
}

// Respond without doing anything.
// This endpoint is used to check that the service is up.
//
// See https://docs.taskcluster.net/reference/core/index/api-docs#ping
func (myIndex *Index) Ping() error {
	cd := tcclient.Client(*myIndex)
	_, _, err := (&cd).APICall(nil, "GET", "/ping", nil, nil)
	return err
}
