// import TestAPI1 from "./api/test_api_1";

export default function Home() {
  async function TestAPI1() {
    const url = 'http://localhost:8000/test?'
    const res = await fetch('http://localhost:8000/test?x=8', {
        method: 'GET',
        // mode: 'no-cors',
        headers: {
            'Access-Control-Allow-Origin': 'http://localhost:8000/',
        },
    });
    console.log(res)
    
    const body = await res.json()
    console.log(body)
  }
  // TestAPI1()
  return (
    <main>
      <h1>hello nextjs</h1>
    </main>
  )
}