package ts3

type Client struct {
	CID              string `json:"cid"`
	CLID             string `json:"clid"`
	ClientDatabaseId string `json:"client_database_id"`
	ClientNickname   string `json:"client_nickname"`
	ClientType       string `json:"client_type"`
}

type VServerInfo struct {
	ConnectionBandwidthReceivedLastMinuteTotal          string `json:"connection_bandwidth_received_last_minute_total"`
	ConnectionBandwidthReceivedLastSecondTotal          string `json:"connection_bandwidth_received_last_second_total"`
	ConnectionBandwidthSentLastMinuteTotal              string `json:"connection_bandwidth_sent_last_minute_total"`
	ConnectionBandwidthSentLastSecondTotal              string `json:"connection_bandwidth_sent_last_second_total"`
	ConnectionBytesReceivedControl                      string `json:"connection_bytes_received_control"`
	ConnectionBytesReceivedKeepalive                    string `json:"connection_bytes_received_keepalive"`
	ConnectionBytesReceivedSpeech                       string `json:"connection_bytes_received_speech"`
	ConnectionBytesReceivedTotal                        string `json:"connection_bytes_received_total"`
	ConnectionBytesSentControl                          string `json:"connection_bytes_sent_control"`
	ConnectionBytesSentKeepalive                        string `json:"connection_bytes_sent_keepalive"`
	ConnectionBytesSentSpeech                           string `json:"connection_bytes_sent_speech"`
	ConnectionBytesSentTotal                            string `json:"connection_bytes_sent_total"`
	ConnectionFiletransferBandwidthReceived             string `json:"connection_filetransfer_bandwidth_received"`
	ConnectionFiletransferBandwidthSent                 string `json:"connection_filetransfer_bandwidth_sent"`
	ConnectionFiletransferBytesReceivedTotal            string `json:"connection_filetransfer_bytes_received_total"`
	ConnectionFiletransferBytesSentTotal                string `json:"connection_filetransfer_bytes_sent_total"`
	ConnectionPacketsReceivedControl                    string `json:"connection_packets_received_control"`
	ConnectionPacketsReceivedKeepalive                  string `json:"connection_packets_received_keepalive"`
	ConnectionPacketsReceivedSpeech                     string `json:"connection_packets_received_speech"`
	ConnectionPacketsReceivedTotal                      string `json:"connection_packets_received_total"`
	ConnectionPacketsSentControl                        string `json:"connection_packets_sent_control"`
	ConnectionPacketsSentKeepalive                      string `json:"connection_packets_sent_keepalive"`
	ConnectionPacketsSentSpeech                         string `json:"connection_packets_sent_speech"`
	ConnectionPacketsSentTotal                          string `json:"connection_packets_sent_total"`
	VirtualserverAntifloodPointsNeededCommandBlock      string `json:"virtualserver_antiflood_points_needed_command_block"`
	VirtualserverAntifloodPointsNeededIpBlock           string `json:"virtualserver_antiflood_points_needed_ip_block"`
	VirtualserverAntifloodPointsNeededPluginBlock       string `json:"virtualserver_antiflood_points_needed_plugin_block"`
	VirtualserverAntifloodPointsTickReduce              string `json:"virtualserver_antiflood_points_tick_reduce"`
	VirtualserverAskForPrivilegekey                     string `json:"virtualserver_ask_for_privilegekey"`
	VirtualserverAutostart                              string `json:"virtualserver_autostart"`
	VirtualserverCapabilityExtensions                   string `json:"virtualserver_capability_extensions"`
	VirtualserverChannelTempDeleteDelayDefault          string `json:"virtualserver_channel_temp_delete_delay_default"`
	VirtualserverChannelsonline                         string `json:"virtualserver_channelsonline"`
	VirtualserverClientConnections                      string `json:"virtualserver_client_connections"`
	VirtualserverClientsonline                          string `json:"virtualserver_clientsonline"`
	VirtualserverCodecEncryptionMode                    string `json:"virtualserver_codec_encryption_mode"`
	VirtualserverComplainAutobanCount                   string `json:"virtualserver_complain_autoban_count"`
	VirtualserverComplainAutobanTime                    string `json:"virtualserver_complain_autoban_time"`
	VirtualserverComplainRemoveTime                     string `json:"virtualserver_complain_remove_time"`
	VirtualserverCreated                                string `json:"virtualserver_created"`
	VirtualserverDefaultChannelAdminGroup               string `json:"virtualserver_default_channel_admin_group"`
	VirtualserverDefaultChannelGroup                    string `json:"virtualserver_default_channel_group"`
	VirtualserverDefaultServerGroup                     string `json:"virtualserver_default_server_group"`
	VirtualserverDownloadQuota                          string `json:"virtualserver_download_quota"`
	VirtualserverFileStorageClass                       string `json:"virtualserver_file_storage_class"`
	VirtualserverFilebase                               string `json:"virtualserver_filebase"`
	VirtualserverFlagPassword                           string `json:"virtualserver_flag_password"`
	VirtualserverHostbannerGfxInterval                  string `json:"virtualserver_hostbanner_gfx_interval"`
	VirtualserverHostbannerGfxUrl                       string `json:"virtualserver_hostbanner_gfx_url"`
	VirtualserverHostbannerMode                         string `json:"virtualserver_hostbanner_mode"`
	VirtualserverHostbannerUrl                          string `json:"virtualserver_hostbanner_url"`
	VirtualserverHostbuttonGfxUrl                       string `json:"virtualserver_hostbutton_gfx_url"`
	VirtualserverHostbuttonTooltip                      string `json:"virtualserver_hostbutton_tooltip"`
	VirtualserverHostbuttonUrl                          string `json:"virtualserver_hostbutton_url"`
	VirtualserverHostmessage                            string `json:"virtualserver_hostmessage"`
	VirtualserverHostmessageMode                        string `json:"virtualserver_hostmessage_mode"`
	VirtualserverIconId                                 string `json:"virtualserver_icon_id"`
	VirtualserverId                                     string `json:"virtualserver_id"`
	VirtualserverIp                                     string `json:"virtualserver_ip"`
	VirtualserverLogChannel                             string `json:"virtualserver_log_channel"`
	VirtualserverLogClient                              string `json:"virtualserver_log_client"`
	VirtualserverLogFiletransfer                        string `json:"virtualserver_log_filetransfer"`
	VirtualserverLogPermissions                         string `json:"virtualserver_log_permissions"`
	VirtualserverLogQuery                               string `json:"virtualserver_log_query"`
	VirtualserverLogServer                              string `json:"virtualserver_log_server"`
	VirtualserverMachineId                              string `json:"virtualserver_machine_id"`
	VirtualserverMaxDownloadTotalBandwidth              string `json:"virtualserver_max_download_total_bandwidth"`
	VirtualserverMaxUploadTotalBandwidth                string `json:"virtualserver_max_upload_total_bandwidth"`
	VirtualserverMaxclients                             string `json:"virtualserver_maxclients"`
	VirtualserverMinAndroidVersion                      string `json:"virtualserver_min_android_version"`
	VirtualserverMinClientVersion                       string `json:"virtualserver_min_client_version"`
	VirtualserverMinClientsInChannelBeforeForcedSilence string `json:"virtualserver_min_clients_in_channel_before_forced_silence"`
	VirtualserverMinIosVersion                          string `json:"virtualserver_min_ios_version"`
	VirtualserverMonthBytesDownloaded                   string `json:"virtualserver_month_bytes_downloaded"`
	VirtualserverMonthBytesUploaded                     string `json:"virtualserver_month_bytes_uploaded"`
	VirtualserverName                                   string `json:"virtualserver_name"`
	VirtualserverNamePhonetic                           string `json:"virtualserver_name_phonetic"`
	VirtualserverNeededIdentitySecurityLevel            string `json:"virtualserver_needed_identity_security_level"`
	VirtualserverNickname                               string `json:"virtualserver_nickname"`
	VirtualserverPassword                               string `json:"virtualserver_password"`
	VirtualserverPlatform                               string `json:"virtualserver_platform"`
	VirtualserverPort                                   string `json:"virtualserver_port"`
	VirtualserverPrioritySpeakerDimmModificator         string `json:"virtualserver_priority_speaker_dimm_modificator"`
	VirtualserverQueryClientConnections                 string `json:"virtualserver_query_client_connections"`
	VirtualserverQueryclientsonline                     string `json:"virtualserver_queryclientsonline"`
	VirtualserverReservedSlots                          string `json:"virtualserver_reserved_slots"`
	VirtualserverStatus                                 string `json:"virtualserver_status"`
	VirtualserverTotalBytesDownloaded                   string `json:"virtualserver_total_bytes_downloaded"`
	VirtualserverTotalBytesUploaded                     string `json:"virtualserver_total_bytes_uploaded"`
	VirtualserverTotalPacketlossControl                 string `json:"virtualserver_total_packetloss_control"`
	VirtualserverTotalPacketlossKeepalive               string `json:"virtualserver_total_packetloss_keepalive"`
	VirtualserverTotalPacketlossSpeech                  string `json:"virtualserver_total_packetloss_speech"`
	VirtualserverTotalPacketlossTotal                   string `json:"virtualserver_total_packetloss_total"`
	VirtualserverTotalPing                              string `json:"virtualserver_total_ping"`
	VirtualserverUniqueIdentifier                       string `json:"virtualserver_unique_identifier"`
	VirtualserverUploadQuota                            string `json:"virtualserver_upload_quota"`
	VirtualserverUptime                                 string `json:"virtualserver_uptime"`
	VirtualserverVersion                                string `json:"virtualserver_version"`
	VirtualserverWeblistEnabled                         string `json:"virtualserver_weblist_enabled"`
	VirtualserverWelcomemessage                         string `json:"virtualserver_welcomemessage"`
}

type Status struct {
	Code             int    `json:"code"`
	FailedPermission string `json:"failed_permission,omitempty"`
	ExtraMessage     string `json:"extra_message,omitempty"`
	Message          string `json:"message"`
}

type Envelope[T any] struct {
	Body   []T    `json:"body,omitempty"` // this only sometimes
	Status Status `json:"status"`         // this always appears
}
