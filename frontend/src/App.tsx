import { useCallback, useState } from 'react'
import './App.css'

const baseUrl = '/api/v1'

const saveBlobImage = (blob: Blob, filename: string) => {
  const link = document.createElement("a");
  link.href = URL.createObjectURL(blob);
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(link.href); // Clean up memory
};


const downloadFile = async (id: number) => {
  try {
    // You can write the URL of your server or any other endpoint used for file upload
    const result = await fetch(`${baseUrl}/download`, {
      method: 'POST',
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ id }),
    });

    const data = await result.blob();
    saveBlobImage(data, "file.dcm");

    console.log(data);
  } catch (error) {
    console.error(error);
  }
}


const ListComponent: React.FC<{values: Array<Record<string, any>>}> = ({ values }) => {
  if (values.length === 0) {
    return null
  }

  return (
    <ul>
      {values.map((value) => {
        const func = () => downloadFile(value.id);

        return (
          <li key={value.id}>
            <p>{value.original_file} {value.series_description}</p>
            <button onClick={func}>Download</button>
          </li>
        )
      })}
    </ul>
  )
}


function App() {
  const [file, setFile] = useState<File | null>(null);
  const [list, setList] = useState<Array<Record<string, any>>>([])

  const handleFileChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  }, []);

  const sendFile = useCallback(async () => {
    if (file) {
      const formData = new FormData();
      formData.append("file", file);

      try {
        // You can write the URL of your server or any other endpoint used for file upload
        const result = await fetch(`${baseUrl}/upload`, {
          method: 'POST',
          body: formData,
        });

        const data = await result.json();

        console.log(data);
      } catch (error) {
        console.error(error);
      }
    }
  }, [file]);

  const listFiles = useCallback(async () => {
    try {
      const response = await fetch(`${baseUrl}/list`, {
        method: 'POST',
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ offset: 0, limit: 100 })
      });

      const res = await response.json();
      setList(res)
    } catch (e) {
      console.log("ERROR", e);
    }
  }, []);

  return (
    <>
      <input id="file" type="file" onChange={handleFileChange} />
      <button onClick={sendFile}>Upload</button>
      <button onClick={listFiles}>Get file list</button>
      <ListComponent values={list} />
    </>
  )
}

export default App
