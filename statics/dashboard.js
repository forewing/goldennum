Vue.component('dashboard', {
    props: [],
    template: "#dashboard-component",
    data: function () {
        return {
            roomHistoryCtx: null,
            roomHistoryChart: null,
            roomHistoryData: {
                type: "line",
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Goldennum',
                        backgroundColor: 'rgba(151, 216, 178, 0.2)',
                        borderColor: 'rgba(151, 216, 178, 1)',
                        pointHitRadius: 10,
                        data: [],
                        fill: false,
                    }]
                },
                options: {
                    legend: false,
                    responsive: true,
                    title: {
                        display: true,
                        text: 'Number History'
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
            userRankData: {
                type: 'horizontalBar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Scores',
                        backgroundColor: 'rgba(120, 213, 215, 0.2)',
                        borderColor: 'rgba(120, 213, 215, 0.8)',
                        data: [],
                        fill: false,
                    }]
                },
                options: {
                    elements: {
                        rectangle: {
                            borderWidth: 1,
                        }
                    },
                    maintainAspectRatio: false,
                    responsive: true,
                    legend: false,
                    title: {
                        display: true,
                        text: 'Ranks'
                    },
                    scales: {
                        xAxes: [{
                            display: true,
                            scaleLabel: {
                                display: true,
                                labelString: 'Score'
                            },
                            ticks: {
                                suggestedMin: 0,
                                suggestedMax: 20,
                            }
                        }],
                        yAxes: [{
                            display: true,
                            scaleLabel: {
                                display: false,
                                labelString: 'User'
                            },
                        }]
                    }
                },
            },
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
                this.refreshRoom(data.RoomHistorys);
                this.refreshUser(data.Users);
            }).catch(error => {
                console.error(error);
            });
        },
        refreshRoom(data) {
            if (!data) {
                return;
            }
            this.roomHistoryChart.data.labels = [];
            this.roomHistoryChart.data.datasets[0].data = [];
            data.sort((a, b) => a.Round - b.Round);
            for (const history of data) {
                if (history.Round == null || history.GoldenNum == null) {
                    continue;
                }
                this.roomHistoryChart.data.labels.push(history.Round);
                this.roomHistoryChart.data.datasets[0].data.push(history.GoldenNum);
            }
            this.roomHistoryChart.update();
        },
        refreshUser(data) {
            if (!data) {
                return;
            }
            this.userRankChart.data.labels = [];
            this.userRankChart.data.datasets[0].data = [];
            data.sort((a, b) => b.Score - a.Score);
            for (const user of data) {
                if (user.Name == null || user.Score == null) {
                    continue;
                }
                this.userRankChart.data.labels.push(user.Name);
                this.userRankChart.data.datasets[0].data.push(user.Score);
            }
            this.updateUserChartSize();
            this.userRankChart.update();
        },
        updateUserChartSize() {
            let len = this.userRankChart.data.labels.length;
            if (len <= 0) {
                len = 1;
            }
            this.$refs.userRankDiv.style.height = `${100 + len * 50}px`;
        },
    },
    mounted() {
        this.roomHistoryCtx = this.$refs.roomHistory.getContext('2d')
        this.roomHistoryChart = new Chart(this.roomHistoryCtx, this.roomHistoryData);
        this.userRankCtx = this.$refs.userRank.getContext('2d');
        this.userRankChart = new Chart(this.userRankCtx, this.userRankData);
        this.updateUserChartSize();
    },
})
