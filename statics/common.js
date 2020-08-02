const KEY_SAVED_ROOM_ID = "SAVED_ROOM_ID";
const KEY_SAVED_USER_ID = "SAVED_USER_ID";
const KEY_SAVED_PASSWORD = "SAVED_PASSWORD";
const KEY_SAVED_SIGNED_IN = "SAVED_SIGNED_IN";

function getSavedRoomId() {
    let savedRoomId = parseInt(localStorage.getItem(KEY_SAVED_ROOM_ID));
    if (savedRoomId > 0) {
        return savedRoomId;
    }
    console.error(`getSavedRoomId invalid: ${savedRoomId}`);
    return 1;
}

function setSavedRoomId(id) {
    if (parseInt(id) > 0) {
        localStorage.setItem(KEY_SAVED_ROOM_ID, id);
        return true;
    }
    console.error(`setSavedRoomId invalid: ${id}`)
    return false;
}

const URL_ROOM_LIST = BASE_URL + "/rooms/"
const URL_ROOM_INFO = BASE_URL + "/room/"
const URL_ROOM_SYNC = BASE_URL + "/sync/"
const URL_USER_CREATE = BASE_URL + "/users/"
const URL_USER_INFO = BASE_URL + "/user/"
const URL_USER_SUBMIT = BASE_URL + "/user/"
const URL_USER_AUTH = BASE_URL + "/user/"

function jsonResponseHandler(response) {
    if (!response.ok) {
        throw {
            error: response.statusText,
            data: response.json(),
        }
    }
    return response.json()
}

async function fetchPostData(url, data) {
    const response = await fetch(url, {
        body: JSON.stringify(data),
        method: 'POST',
        headers: new Headers({
            'Content-Type': 'application/json'
        }),
    });
    return jsonResponseHandler(response);
}

async function fetchPutData(url, data) {
    const response = await fetch(url, {
        body: JSON.stringify(data),
        method: 'PUT',
        headers: new Headers({
            'Content-Type': 'application/json'
        }),
    });
    return jsonResponseHandler(response);
}

async function fetchGetData(url) {
    const response = await fetch(url, {
        method: 'GET',
    });
    return jsonResponseHandler(response);
}

function getRoomList() {
    return fetchGetData(URL_ROOM_LIST)
}

function getRoomInfo(roomId) {
    return fetchGetData(URL_ROOM_INFO + parseInt(roomId))
}

function getRoomSync(roomId) {
    return fetchGetData(URL_ROOM_SYNC + parseInt(roomId))
}

function postUserCreate(roomId, username, password) {
    return fetchPostData(URL_USER_CREATE + parseInt(roomId), {
        Username: String(username),
        Password: String(password),
    })
}

function getUserInfo(userId) {
    return fetchGetData(URL_USER_INFO + parseInt(userId))
}

function postUserSubmit(userId, password, submit1, submit2) {
    return fetchPostData(URL_USER_SUBMIT + parseInt(userId), {
        Password: String(password),
        Submit1: parseFloat(submit1),
        Submit2: parseFloat(submit2),
    })
}

function putUserAuth(userId, password) {
    return fetchPutData(URL_USER_AUTH + parseInt(userId), {
        Password: String(password)
    })
}
