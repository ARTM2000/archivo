package sourceserver

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const defaultBucketDurationSize = 5 * time.Second // default is 5 seconds
const defaultBucketStoreSize = 3 * 24 * time.Hour // default is 3 days

var buckets []Bucket = []Bucket{}
var mu sync.Mutex

func init() {
	createEmptyBucket(defaultBucketDurationSize)
	go func() {
		ticker := time.NewTicker(defaultBucketDurationSize)
		for {
			<-ticker.C
			createEmptyBucket(defaultBucketDurationSize)
		}
	}()

	go func() {
		storeTimer := time.NewTimer(defaultBucketStoreSize)
		storeTicker := time.NewTicker(defaultBucketDurationSize)
		<-storeTimer.C
		log.Default().Println("start rotating....")
		for {
			<-storeTicker.C
			rotateStoredBuckets(defaultBucketStoreSize)
		}
	}()
}

const (
	FailOperation = iota
	SuccessOperation
)

type BucketDetail struct {
	SuccessCount int64 `json:"success_count"`
	FailCount    int64 `json:"fail_count"`
}

type Bucket struct {
	From         time.Time
	To           time.Time
	TotalSuccess int64
	TotalFail    int64
	SrvAddress   map[string]int
	Details      []BucketDetail // holds each source server metrics by key of source server
}

type BucketReport struct {
	From         time.Time               `json:"from"`
	To           time.Time               `json:"to"`
	TotalSuccess int64                   `json:"total_success"`
	TotalFail    int64                   `json:"total_fail"`
	Details      map[string]BucketDetail `json:"details"`
}

func NewSrcSrvMetrics() *srcSrvMetrics {
	return &srcSrvMetrics{}
}

type srcSrvMetrics struct{}

func createEmptyBucket(timeSpan time.Duration) {
	if timeSpan == 0 {
		timeSpan = defaultBucketDurationSize
	}
	now := time.Now()

	mu.Lock()
	defer mu.Unlock()

	initialDetail := []BucketDetail{}
	initialSrvAddress := map[string]int{}
	bucket := Bucket{
		From:       now,
		To:         now.Add(defaultBucketDurationSize),
		Details:    initialDetail,
		SrvAddress: initialSrvAddress,
	}
	buckets = append(buckets, bucket)
}

func rotateStoredBuckets(rotateTime time.Duration) {
	if rotateTime == 0 {
		rotateTime = defaultBucketStoreSize
	}

	lastTime := time.Now().Add(-1 * rotateTime)

	mu.Lock()
	defer mu.Unlock()

	for {
		i := 0
		if lastTime.Before(buckets[i].To) {
			break
		}
		buckets = buckets[i+1:]
	}
}

func (ssm *srcSrvMetrics) filterBucketsByTime(from, to time.Time) (b []Bucket) {
	log.Default().Println("from:", from, "to:", to)
	for _, bc := range buckets {
		if bc.To.Before(to) && bc.From.After(from) {
			b = append(b, bc)
		}
	}
	return
}

func (ssm *srcSrvMetrics) formatBuckets(bcs []Bucket) []BucketReport {
	bcr := []BucketReport{}
	for _, bc := range bcs {
		br := BucketReport{
			From:         bc.From,
			To:           bc.To,
			TotalSuccess: bc.TotalSuccess,
			TotalFail:    bc.TotalFail,
			Details:      map[string]BucketDetail{},
		}
		for k, v := range bc.SrvAddress {
			br.Details[k] = bc.Details[v]
		}
		bcr = append(bcr, br)
	}
	return bcr
}

func (ssm *srcSrvMetrics) CountOperation(sourceServerName string, status int) {
	mu.Lock()
	defer mu.Unlock()
	bucketLength := len(buckets)
	lastBucket := buckets[bucketLength-1]

	// in case that operation can place in last bucket...
	index, exists := lastBucket.SrvAddress[sourceServerName]
	if exists {
		if status == SuccessOperation {
			atomic.AddInt64(&buckets[len(buckets)-1].Details[index].SuccessCount, 1)
			atomic.AddInt64(&buckets[len(buckets)-1].TotalSuccess, 1)
		} else {
			atomic.AddInt64(&buckets[len(buckets)-1].Details[index].FailCount, 1)
			atomic.AddInt64(&buckets[len(buckets)-1].TotalFail, 1)
		}
		return
	}

	if status == SuccessOperation {
		buckets[bucketLength-1].Details = append(buckets[bucketLength-1].Details, BucketDetail{SuccessCount: 1})
		atomic.AddInt64(&buckets[bucketLength-1].TotalSuccess, 1)
	} else {
		buckets[bucketLength-1].Details = append(buckets[bucketLength-1].Details, BucketDetail{FailCount: 1})
		atomic.AddInt64(&buckets[bucketLength-1].TotalFail, 1)
	}
	buckets[bucketLength-1].SrvAddress[sourceServerName] = len(lastBucket.Details)
}

func (ssm *srcSrvMetrics) SingleSrvBucketsAsMetrics(sourceServerName string, from, to time.Time) []BucketReport {
	buckets := ssm.filterBucketsByTime(from, to)
	for i, b := range buckets {
		index := i
		sAddr, exists := b.SrvAddress[sourceServerName]
		if !exists {
			// if source server not exists in bucket, clear any unrelated data
			buckets[index].TotalSuccess = 0
			buckets[index].TotalFail = 0
			buckets[index].SrvAddress = nil
			buckets[index].Details = []BucketDetail{}
			continue
		}

		sBucket := buckets[index].Details[sAddr]
		buckets[index].TotalSuccess = sBucket.SuccessCount
		buckets[index].TotalFail = sBucket.FailCount
		buckets[index].SrvAddress = map[string]int{
			sourceServerName: 0,
		}
		buckets[index].Details = []BucketDetail{
			sBucket,
		}
	}
	return ssm.formatBuckets(buckets)
}

func (ssm *srcSrvMetrics) AllBucketsAsMetrics(from, to time.Time) []BucketReport {
	return ssm.formatBuckets(ssm.filterBucketsByTime(from, to))
}
