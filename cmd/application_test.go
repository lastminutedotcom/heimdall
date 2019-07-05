package cmd

import (
	"fmt"
	"github.com/lastminutedotcom/heimdall/pkg/client/colocation"
	"github.com/lastminutedotcom/heimdall/pkg/client/ratelimit"
	"github.com/lastminutedotcom/heimdall/pkg/client/waf"
	"github.com/lastminutedotcom/heimdall/pkg/client/zone"
	"github.com/lastminutedotcom/heimdall/pkg/metric"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/marpaia/graphite-golang"
	"net"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

const sampleZoneID = "aaaaaaaabbbbbbbbccccccc"

func Test_integrationTest(t *testing.T) {

	mockZones := zone.MockZones{
		Path: filepath.Join("..", "test", "IT", "cloudflare_zone.json"),
	}

	mockColocations := colocation.MockColocations{
		Path: filepath.Join("..", "test", "IT", "cloudflare_colocation.json"),
	}

	mockWafs := waf.MockWafs{
		Path: filepath.Join("..", "test", "IT", "cloudflare_waf.json"),
	}

	mockRateLimitClient := ratelimit.MockRateLimitClient{
		Path: filepath.Join("..", "test", "IT", "cloudflare_ratelimit.json"),
	}

	aggregate := collect(&model.Config{}, mockZones, mockColocations, mockWafs, mockRateLimitClient)
	metrics := metric.AdaptDataToMetrics(aggregate)

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Name < metrics[j].Name
	})

	assert.Equal(t, len(metrics), 337)
	assert.Equal(t, metrics[15].Name, "cloudflare.play_at.secure_play_at.total.ratelimit.post.challenge")
	assert.Equal(t, metrics[24].Name, "cloudflare.play_at.secure_play_at.total.ratelimit.put.simulate")
	assert.Equal(t, metrics[15].Value, "4")
	assert.Equal(t, metrics[24].Value, "1")
}

// Test to emulate the streaming behavior
func Test_metricIsReadAndStreamedToGraphite(t *testing.T) {
	// mock a graphite UDP server
	done := make(chan bool)
	a, err := net.ResolveUDPAddr("udp", ":2113")
	if err != nil {
		t.Fatalf("could not listen on UDP addr :2113")
	}
	go mockedUDPserver(t, done, a)
	//shutdown the mocked UDP server when the test is over
	defer func() {
		done <- true
	}()

	g, err := graphite.NewGraphiteUDP("127.0.0.1", 2113)
	if err != nil {
		t.Fatalf("could not connect to graphite")
	}
	err = g.SendMetric(graphite.Metric{Name: "a.sample.metric", Value: "12345", Timestamp: time.Now().Unix()})
	assert.NoError(t, err)

	metrics := []*model.Aggregate{
		{ZoneName: "test.com", ZoneID: sampleZoneID, Totals: smallTotalsSample()},
		{ZoneName: "test.de", ZoneID: sampleZoneID, Totals: smallTotalsSample()},
		{ZoneName: "test.uk", ZoneID: sampleZoneID, Totals: smallTotalsSample()},
		{ZoneName: "test.it", ZoneID: sampleZoneID, Totals: smallTotalsSample()},
	}

	// crucial to have a buffered channel of the exact size of the metrics to be sent!!
	testAggrCh := make(chan *model.Aggregate, len(metrics))
	for _,m := range metrics {
		testAggrCh <- m
	}
	// otherwise the close here will not work and there will be stuck goroutine, thanks Bill Kennedy
	close(testAggrCh)
	if err := adaptAndSend(testAggrCh, g); err!=nil {
		t.Errorf("error converting metrics and pushing to Graphite: %v", err)
	}
}

func smallTotalsSample() map[time.Time]*model.Counters {
	m := make(map[time.Time]*model.Counters, 0)
	m[time.Unix(123456789, 0)] = &model.Counters{
		RequestAll: model.Counter{
			Key:   "2xx",
			Value: 200,
		},
	}
	return m
}

func mockedUDPserver(t *testing.T, done <-chan bool, addr *net.UDPAddr) {
	b := make([]byte, 4096)
	server, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("could not start UDP stub server")
	}
	for {
		select {
		default:
			n, conn, err := server.ReadFromUDP(b)
			if err != nil {
				t.Fatalf("error reading UDP stream: %v", err)
				break
			}
			if conn == nil {
				continue
			}
			// echo the input to the client
			s, err := server.WriteToUDP(b[:n], conn)
			if err != nil {
				t.Fatalf("error responding from UDP server to client: %v", err)
			}
			// ensure we echo back the input
			assert.Equal(t, n, s)

		case <-done:
			fmt.Println("UDP mocked server stopped!")
			break
		}
	}
}
