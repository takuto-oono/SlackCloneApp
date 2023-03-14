//Cookieに保存されているjwtTokenを取り出す
export function getToken() {
  if (typeof document !== 'undefined') {
    try {
      const jwtToken = document.cookie.split('; ').find(row => row.startsWith('token')).split('=')[1];
      return jwtToken;
      
    } catch (err) {
      console.log(err)
      return;
    }
  }
}

//Cookieに保存されているuser_idを取り出す
export function getUserId() {
  if (typeof document !== 'undefined') {
    try {
      const userId = document.cookie.split('; ').find(row => row.startsWith('user_id')).split('=')[1];
      return parseInt(userId);
    } catch (err) {
      console.log(err)
      return;
    }
  }
}
