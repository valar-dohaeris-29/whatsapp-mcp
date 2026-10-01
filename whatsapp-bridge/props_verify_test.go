package main

import (
	"testing"

	"go.mau.fi/whatsmeow/store"
	"google.golang.org/protobuf/proto"
)

// Applies the same mutations main() does, then checks that raising the limits
// did not clobber whatsmeow's default capability flags.
func TestHistorySyncPropsPreserveDefaults(t *testing.T) {
	store.DeviceProps.RequireFullSync = proto.Bool(true)
	store.DeviceProps.HistorySyncConfig.FullSyncDaysLimit = proto.Uint32(3650)
	store.DeviceProps.HistorySyncConfig.FullSyncSizeMbLimit = proto.Uint32(2048)
	store.DeviceProps.HistorySyncConfig.StorageQuotaMb = proto.Uint32(10240)

	hsc := store.DeviceProps.HistorySyncConfig

	if !store.DeviceProps.GetRequireFullSync() {
		t.Error("RequireFullSync not set")
	}
	t.Logf("RequireFullSync     = %v", store.DeviceProps.GetRequireFullSync())
	t.Logf("FullSyncDaysLimit   = %d", hsc.GetFullSyncDaysLimit())
	t.Logf("FullSyncSizeMbLimit = %d", hsc.GetFullSyncSizeMbLimit())
	t.Logf("StorageQuotaMb      = %d", hsc.GetStorageQuotaMb())

	// These come from whatsmeow's defaults and must survive.
	defaults := map[string]bool{
		"SupportGroupHistory":           hsc.GetSupportGroupHistory(),
		"SupportCallLogHistory":         hsc.GetSupportCallLogHistory(),
		"SupportCagReactionsAndPolls":   hsc.GetSupportCagReactionsAndPolls(),
		"InlineInitialPayloadInE2EeMsg": hsc.GetInlineInitialPayloadInE2EeMsg(),
		"SupportBizHostedMsg":           hsc.GetSupportBizHostedMsg(),
		"SupportMessageAssociation":     hsc.GetSupportMessageAssociation(),
	}
	for name, got := range defaults {
		if !got {
			t.Errorf("default flag %s was clobbered (got false)", name)
		}
		t.Logf("%-30s = %v", name, got)
	}

	// The payload must actually serialize.
	b, err := proto.Marshal(store.DeviceProps)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	t.Logf("serialized DeviceProps = %d bytes", len(b))
}
