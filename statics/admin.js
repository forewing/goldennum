var app = new Vue({
    el: '#app',
    data: {
        rooms: [],
        inputRoomInterval: 20,
        inputRoomTotalRounds: 30,
    },
    methods: {
        refreshRooms() {
            getRoomList().then(data => {
                this.rooms = [];
                for (const room of data) {
                    this.rooms.unshift(room.ID);
                }
            }).catch(error => {
                console.log(error);
                if (error.data) {
                    error.data.then(data => console.log(data));
                }
            })
        },
        addRoom() {
            postCreateRoom(this.inputRoomInterval, this.inputRoomTotalRounds).then(data => {
                this.rooms.unshift(data.ID);
            }).catch(error => console.log(error));
        },
    },
    mounted: function () {
        this.refreshRooms();
    },
})