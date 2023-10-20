package sourceserver

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const defaultBucketDurationSize = 5 * time.Second

var buckets []Bucket = []Bucket{}
var mu sync.Mutex

func init() {
	createEmptyBucket()
	go func() {
		ticker := time.NewTicker(defaultBucketDurationSize)
		for {
			<-ticker.C
			createEmptyBucket()
		}
	}()
}

const (
	FailOperation = iota
	SuccessOperation
)

type BucketDetail struct {
	SuccessCount int64
	FailCount    int64
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

func createEmptyBucket() {
	now := time.Now()
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
			From: bc.From,
			To: bc.To,
			TotalSuccess: bc.TotalSuccess,
			TotalFail: bc.TotalFail,
			Details: map[string]BucketDetail{},
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
		buckets[bucketLength-1].TotalSuccess = 1
	} else {
		buckets[bucketLength-1].Details = append(buckets[bucketLength-1].Details, BucketDetail{FailCount: 1})
		buckets[bucketLength-1].TotalFail = 1
	}
	buckets[bucketLength-1].SrvAddress[sourceServerName] = len(lastBucket.Details)
}

func (ssm *srcSrvMetrics) SingleSrvBucketsAsMetrics(sourceServerName string, from, to time.Time) (srvBuckets []Bucket) {
	for _, b := range ssm.filterBucketsByTime(from, to) {
		sAddrs, exist := b.SrvAddress[sourceServerName]
		if !exist {
			continue
		}
		sBucket := b.Details[sAddrs]
		srb := Bucket{
			From:         b.From,
			To:           b.To,
			TotalSuccess: sBucket.SuccessCount,
			TotalFail:    sBucket.FailCount,
		}
		srvBuckets = append(srvBuckets, srb)
	}

	return srvBuckets
}

func (ssm *srcSrvMetrics) AllBucketsAsMetrics(from, to time.Time) []BucketReport {
	return ssm.formatBuckets(ssm.filterBucketsByTime(from, to))
}
