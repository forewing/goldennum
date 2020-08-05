var app = new Vue({
    el: '#app',
    data: {
        rooms: [],
    },
    methods: {
        refreshRooms() {
            getRoomList().then(data => {
                this.rooms = [];
                for (const room of data) {
                    this.rooms.push(room.ID);
                }
            }).catch(error => {
                console.log(error);
                if (error.data) {
                    error.data.then(data => console.log(data));
                }
            })
        },
    },
    mounted: function () {
        this.refreshRooms();
    },
})