Vue.component('dashboard', {
    delimiters: ['[[', ']]'],
    props: [],
    template: "#dashboard-component",
    data: function () {
        return {
            data: null,
            roomId: 1,
            intervalId: null,
            isButtonStart: true,
            nextTick: Date.now(),
            countDown: 0,
            errorMessage: "",
            historyLength: 0,
            roomHistoryCtx: null,
            roomHistoryChart: null,
            roomHistoryData: {
                type: "line",
                data: {
                    labels: [],
                    datasets: [{
                        label: DATA_LABEL_GOLDENNUM,
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
                        text: TITLE_NUMBER_HISTORY
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
                                labelString: SCALE_LABEL_ROUNDS
                            }
                        }],
                        yAxes: [{
                            display: true,
                            scaleLabel: {
                                display: true,
                                labelString: SCALE_LABEL_GOLDENNUM
                            },
                            ticks: {
                                suggestedMin: 0,
                                suggestedMax: 20,
                            }
                        }]
                    }
                },
            },
            users: [],
            userHistory: [],
            userHistoryUser: 0,
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
            this.isButtonStart = false;
            this.intervalId = setTimeout(func, timeout);
            this.nextTick = Date.now() + timeout;
        },
        clearTimeout() {
            this.isButtonStart = true;
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
                const seconds = parseFloat(data);
                if (seconds < 0) {
                    this.setTimeout(this.syncRefresh, 5000);
                } else {
                    this.setTimeout(this.syncRefresh, (seconds + 1) * 1000);
                }
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
                this.errorMessage = "";
                this.data = data;
                this.saveUserScores(data.Users);
                this.refreshRoom(data.RoomHistorys);
                this.refreshUser(data.Users);
            }).catch(error => {
                this.errorMessage = error.error;
                if (error.data) {
                    error.data.then(data => this.errorMessage += data.length > 0 ? ", " + data : "")
                }
                console.error(error);
            });
        },
        saveUserScores(data) {
            for (const user of data) {
                if (user.ID == null || user.Score == null) {
                    continue;
                }
                localStorage.setItem(KEY_USER_SCORE_PREFIX + user.ID, user.Score);
            }
        },
        refreshRoom(data) {
            if (!data) {
                data = [];
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
                data = [];
            }
            data = data.filter(user => user.ID != null && user.Score != null && user.Name != null);
            data.sort((a, b) => b.Score - a.Score);
            this.users = data;
        },
        showUserHistory(event) {
            if (!event) {
                return;
            }
            const userId = event.target.attributes.userid.value;
            getUserInfo(userId).then(data => {
                if (!data.UserHistorys) {
                    data.UserHistorys = [];
                }
                this.userHistoryUser = userId;
                this.userHistory = data.UserHistorys.reverse();
                $("#userHistoryModal").modal('show');
            }).catch(error => {
                console.error(error);
                if (error.data) {
                    error.data.then(data => console.error(data));
                }
            });
        },
    },
    mounted() {
        this.roomHistoryCtx = this.$refs.roomHistory.getContext('2d')
        this.roomHistoryChart = new Chart(this.roomHistoryCtx, this.roomHistoryData);
        this.roomId = getSavedRoomId();
        this.toggleRefresh();
        this.refreshTimeOut();
    },
})
