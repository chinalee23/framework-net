using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Collections.Generic;

namespace Net {
    class TransferTcp : Transfer {
        public TransferTcp() {
            msgQueue = new List<byte[]>();
        }
        
        List<byte[]> msgQueue;

        byte[] pack(byte[] data) {
            int len = data.Length;

            byte[] btLen;
            if (len > 127) {
                btLen = new byte[2];
                btLen[0] = (byte)(((len >> 8) & 0x7f) | 0x80);
                btLen[1] = (byte)(len & 0xff);
            } else {
                btLen = new byte[1];
                btLen[0] = (byte)len;
            }
            byte[] btMsg = new byte[btLen.Length + len];
            Array.Copy(btLen, 0, btMsg, 0, btLen.Length);
            Array.Copy(data, 0, btMsg, btLen.Length, len);

            return btMsg;
        }

        bool unpack(ref int offset) {
            int sz;
            int szLen;
            if (dataBuffer.bt[offset] > 127) {
                if (offset == dataBuffer.len) {
                    return false;
                }
                sz = (dataBuffer.bt[offset] & 0x7f) * 256 + dataBuffer.bt[offset + 1];
                szLen = 2;
            } else {
                sz = dataBuffer.bt[offset];
                szLen = 1;
            }
            if ((dataBuffer.len - offset - szLen) < sz) {
                return false;
            }
            byte[] bt = new byte[sz];
            Array.Copy(dataBuffer.bt, offset + szLen, bt, 0, sz);
            msgQueue.Add(bt);

            offset += sz + szLen;
            return true;
        }

        void read() {
            while (true) {
                try {
                    int sz = socket.Receive(socketBuffer.bt, socketBuffer.len, Definition.BUFFER_SIZE - socketBuffer.len, SocketFlags.None);
                    socketBuffer.len += sz;
                    copyToDataBuffer();
                } catch (Exception e) {
                    error(e);
                    return;
                }
            }
        }

        void onConnect(IAsyncResult ar) {
            try {
                socket.EndConnect(ar);
                threadRead = new Thread(new ThreadStart(read));
                threadRead.Start();

                Connected = true;
            }
            catch (Exception e) {
                error(e);
                return;
            }
        }

        public void Connect(IPEndPoint remote) {
            Connected = false;
            try {
                socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
                socket.BeginConnect(remote, onConnect, null);
            } catch (Exception e) {
                error(e);
            }
        }

        public void Send(byte[] data) {
            byte[] msg = pack(data);
            try {
                socket.Send(msg, msg.Length, SocketFlags.None);
            } catch (Exception e) {

            }
        }

        public void Update() {
            lock (dataBuffer) {
                int offset = 0;
                while (unpack(ref offset)) { }
                Array.Copy(dataBuffer.bt, offset, dataBuffer.bt, 0, dataBuffer.len - offset);
            }
        }

        public byte[] Recv() {
            if (msgQueue.Count > 0) {
                byte[] data = msgQueue[0];
                msgQueue.RemoveAt(0);
                return data;
            } else {
                return null;
            }
        }
    }
}
