import React, { useCallback, useEffect, useState } from 'react';

const socket = new WebSocket('ws://127.0.0.1:3000/ws');

function App() {
    const [message, setMessage] = useState('');
    const [inputValue, setInputValue] = useState('');
    const [data, setData] = useState([]);
    const [showPopup, setShowPopup] = useState(false);

    useEffect(() => {
        socket.onopen = () => {
            setMessage('Connected');
        };

        socket.onmessage = (e) => {
            const response = JSON.parse(e.data);
            setData(response.data);
        };

        socket.onclose = () => {
            setMessage('Disconnected');
        };

        try {
            socket.send(
                JSON.stringify({
                    q: inputValue,
                })
            );
        } catch (err) {
            console.error(err);
        }
        return () => {};
    }, [inputValue]);

    const handleChange = useCallback((e) => {
        setInputValue(e.target.value);
    }, []);

    const handleClose = useCallback((e) => {
        e.preventDefault();
        socket.close();
    }, []);

    const togglePopup = () => {
        setShowPopup(!showPopup);
    };

    const highlightKeyword = (text, keyword) => {
        if (keyword === '') return text;
        const regex = new RegExp(`(${keyword})`, 'gi');
        return text.replace(regex, '<span class="bg-sky-300">$1</span>');
    };

    return (
        <main className='bg-black min-h-screen w-full p-12'>
            <div className='relative max-w-lg mx-auto'>
                <input
                    value={inputValue}
                    type='text'
                    className='w-full p-3 rounded-lg shadow-sm border border-gray-500 bg-gray-800 text-gray-100 placeholder-gray-500 focus:border-gray-500 focus:outline-none'
                    placeholder='Search...'
                    onChange={handleChange}
                    onFocus={togglePopup}
                    onBlur={togglePopup}
                />

                <div className='absolute right-0 top-0 mt-3 mr-3'>
                    <button className='text-gray-600 hover:text-gray-700 focus:outline-none'>
                        <span>
                            Terdapat{' '}
                            <span className='text-yellow-500'>
                                {data && data.length}
                            </span>{' '}
                            hasil pencarian
                        </span>
                    </button>
                </div>

                {showPopup && (
                    <div className='bg-white mt-2 rounded-lg shadow-lg absolute z-10 w-full max-h-80 overflow-y-auto'>
                        <ul className='list-reset'>
                            {data &&
                                data.map((item, index) => (
                                    <li key={index}>
                                        <a
                                            href='#'
                                            className='block group px-4 py-2 text-gray-800 hover:bg-indigo-500 hover:text-white'
                                        >
                                            <div
                                                dangerouslySetInnerHTML={{
                                                    __html: highlightKeyword(
                                                        item.name,
                                                        inputValue
                                                    ),
                                                }}
                                            />

                                            <span className='text-xs hidden  group-hover:flex'>
                                                {item.description}
                                            </span>
                                        </a>
                                    </li>
                                ))}
                        </ul>
                    </div>
                )}
            </div>
        </main>
    );
}

export default App;
