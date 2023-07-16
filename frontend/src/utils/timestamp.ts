// backendから送信されてくるデータのフォーマットは「2023-07-15T02:49:03.247Z」みたいな感じ
// Tは日付と、時刻の間に入るだけで大した意味はない、Zはタイムゾーン(イギリス時間)
// メソッドの設計思想としては、backendから送信されてきたフォーマットを引数として与えると、最終的な結果を返せるようにする。表示用のComponentでここにあるようなメソッドを連続で使わないといけなかったりするのはNG
// EX changeFormat(timeStamp).getDate()のような感じにしないと使い物にならないのはNG

export const convertJSTDate = (backendTimeStamp: string): Date => {
  return new Date(backendTimeStamp)
}

const toLocalTimeStringJWT24h = (date: Date): string => {
  return date.toLocaleTimeString('JWT', {
    hour12: false
  })
}

const toLocalDateStringJWTFormat = (date: Date): string => {
  return date.toLocaleDateString('ja-JP', {})
}

export const getTimeStringFromTimeStamp = (timeStamp: string): string => {
  return toLocalTimeStringJWT24h(convertJSTDate(timeStamp))
}

export const getDateStringFromTimeStamp = (timeStamp: string): string => {
  return toLocalDateStringJWTFormat(convertJSTDate(timeStamp))
}

export const isSameDate = (timeStamp1: string, timeStamp2: string): boolean => {
  const d1: Date = convertJSTDate(timeStamp1)
  const d2: Date = convertJSTDate(timeStamp2)
  return d1.getFullYear() == d2.getFullYear() && d1.getMonth() == d2.getMonth() && d1.getDate() == d2.getDate()
}
