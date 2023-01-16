const url = 'http://localhost:8000/1'

export default async function TestAPI1() {
    const res = await fetch(url, {
        method: 'GET',
        
    });
    
    const body = await res.json()
    console.log(body)
}