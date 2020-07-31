Vue.component('dashboard', {
    props: [],
    template: "#dashboard-component",
    data: function () {
        return {
            data: "{}",
            roomHistoryCtx: null,
            roomHistoryChart: null,
            roomHistoryData: {
                type: "line",
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Goldennum',
                        backgroundColor: 'rgb(255, 99, 132)',
                        borderColor: 'rgb(255, 99, 132)',
                        data: [],
                        fill: false,
                    }]
                },
                options: {
                    legend: false,
                    responsive: true,
                    title: {
                        display: true,
                        text: 'Goldennum History'
                    },
                    tooltips: {
                        mode: 'index',
                        intersect: true
                    },
                    scales: {
                        xAxes: [{
                            display: true,
                            scaleLabel: {
                                display: true,
                                labelString: 'Rounds'
                            }
                        }],
                        yAxes: [{
                            display: true,
                            scaleLabel: {
                                display: true,
                                labelString: 'Goldennum'
                            },
                            ticks: {
                                suggestedMin: 0,
                                suggestedMax: 20,
                            }
                        }]
                    }
                },
            },
            userRankCtx: null,
            userRankChart: null,
            userRankData: {},
        }
    },
    methods: {
        refreshRoomAction() {
            localStorage.setItem(KEY_ROOM_ID, 1);
            const roomId = localStorage.getItem(KEY_ROOM_ID);
            if (!roomId) {
                return;
            }
            getRoomInfo(roomId).then(data => {
                this.refreshRoom(data);
            }).catch(error => {
                console.error(error);
            });
        },
        refreshRoom(data) {
            this.data = data;
            if (!data.RoomHistorys) {
                return;
            }
            this.roomHistoryChart.data.labels = [];
            this.roomHistoryChart.data.datasets[0].data = [];
            for (const history of data.RoomHistorys) {
                if (history.Round == null || history.GoldenNum == null) {
                    continue;
                }
                this.roomHistoryChart.data.labels.push(history.Round);
                this.roomHistoryChart.data.datasets[0].data.push(history.GoldenNum);
            }
            this.roomHistoryChart.update();
        }
    },
    mounted() {
        this.roomHistoryCtx = this.$refs.roomHistory.getContext('2d')
        this.roomHistoryChart = new Chart(this.roomHistoryCtx, this.roomHistoryData);
        this.userRankCtx = this.$refs.userRank.getContext('2d');
    },
})
