Vue.component('dashboard', {
    props: [],
    template: "#dashboard-component",
    data: function () {
        return {
            data: null,
            roomId: 1,
            intervalId: null,
            startStopButtonText: "Start",
            nextTick: Date.now(),
            countDown: 0,
            historyLength: 0,
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
                    animation: {
                        duration: 750,
                        easing: 'easeOutBounce',
                    },
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
        updateRoomId(id) {
            this.roomId = id;
            this.refreshWorker();
        },
        updateHistoryLength(len) {
            this.historyLength = len;
            this.refreshRoom(this.data.RoomHistorys);
        },
        refreshTimeOut() {
            if (this.nextTick > Date.now()) {
                this.countDown = parseInt((this.nextTick - Date.now()) / 1000);
            } else {
                this.countDown = 0;
            }
            setTimeout(this.refreshTimeOut, 1000);
        },
        setTimeout(func, timeout) {
            this.startStopButtonText = "Stop";
            this.intervalId = setTimeout(func, timeout);
            this.nextTick = Date.now() + timeout;
        },
        clearTimeout() {
            this.startStopButtonText = "Start";
            clearInterval(this.intervalId);
            this.intervalId = null;
        },
        toggleRefresh() {
            if (this.intervalId == null) {
                this.setTimeout(this.syncRefresh, 100);
                return;
            }
            clearInterval(this.intervalId);
            this.clearTimeout();
        },
        syncRefresh() {
            if (this.intervalId == null) {
                return;
            }
            this.refreshWorker();
            getRoomSync(this.roomId).then(data => {
                const time = parseInt(data) + 1;
                this.setTimeout(this.syncRefresh, time * 1000);
            }).catch(error => {
                this.setTimeout(this.syncRefresh, 5000);
            })
        },
        refreshWorker() {
            const roomId = parseInt(this.roomId);
            if (typeof (roomId) != "number" || roomId <= 0) {
                return;
            }
            getRoomInfo(roomId).then(data => {
                this.data = data;
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
            let leftBound = 0;
            let rightBound = data.length;
            if (this.historyLength > 0) {
                leftBound = rightBound - this.historyLength;
            }
            if (leftBound < 0) {
                leftBound = 0;
            }
            for (const history of data.slice(leftBound, rightBound)) {
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
        let savedRoomId = parseInt(localStorage.getItem(KEY_SAVED_ROOM_ID));
        if (savedRoomId > 0) {
            this.roomId = savedRoomId
        }
        this.updateUserChartSize();
        this.toggleRefresh();
        this.refreshTimeOut();
    },
})
