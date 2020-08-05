Vue.component('room-control', {
    props: ['roomid'],
    template: "#room-control-component",
    data: function () {
        return {
            userNum: 0,
            interval: 0,
            roundNow: 0,
            roundTotal: 0,
            closed: false,
            invalid: false,
            color: "",
        }
    },
    methods: {
        updateRoomInfo() {
            getRoomInfo(this.roomid).then(data => {
                if (Array.isArray(data.Users)) {
                    this.userNum = data.Users.length;
                }
                this.interval = data.Interval;
                this.roundNow = data.RoundNow;
                this.roundTotal = data.RoundTotal;
            }).catch(error => {
                this.invalid = true;
                console.log(error);
                if (error.data) {
                    error.data.then(data => console.log(data));
                }
            })
        },
        updateRoomSync() {
            getRoomSync(this.roomid).then(data => {
                this.closed = false;
            }).catch(error => {
                this.closed = true;
            })
        },
        toggleStatus() {
            if (this.closed) {
                this.startRoom();
                return;
            }
            this.stopRoom();
        },
        startRoom() {
            putStartRoom(this.roomid).then(data => {
                this.blinkSuccess();
                this.updateRoomSync();
            }).catch(error => {
                this.blinkDanger();
            });
        },
        stopRoom() {
            deleteStopRoom(this.roomid).then(data => {
                this.blinkSuccess();
                this.updateRoomSync();
            }).catch(error => {
                this.blinkDanger();
            });
        },
        blinkSuccess() {
            this.color = "table-success";
            setTimeout(() => this.color = "", 1000);
        },
        blinkDanger() {
            this.color = "table-danger";
            setTimeout(() => this.color = "", 1000);
        },
    },
    mounted() {
        this.updateRoomSync();
        this.updateRoomInfo();
    },
})
