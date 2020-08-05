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
    },
    mounted() {
        this.updateRoomSync();
        this.updateRoomInfo();
    },
})
