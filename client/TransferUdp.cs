using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Collections.Generic;

namespace Net {
    class TransferUdp : Transfer {
        public TransferUdp() {
            U = new Rudp();
        }

        EndPoint rep;
        bool binded = false;
        Rudp U;

        void read() {
            while (true) {
                try {
                    if (!binded) {
                        Thread.Sleep(10);
                        continue;
                    }
                    EndPoint ep = new IPEndPoint(IPAddress.Any, 0);
                    int sz = socket.ReceiveFrom(socketBuffer.bt, socketBuffer.len, Definition.BUFFER_SIZE - socketBuffer.len, SocketFlags.None, ref ep);
                    if (!ep.Equals(rep)) {
                        continue;
                    }

                    socketBuffer.len += sz;
                    copyToDataBuffer();
                }
                catch (Exception e) {
                    error(e);
                    return;
                }
            }
        }

        void send(byte[] data, int len) {
            try {
                socket.SendTo(data, data.Length, SocketFlags.None, rep);
                if (!binded) {
                    binded = true;
                }
            }
            catch (Exception e) {
                error(e);
            }
        }

        public void Connect(IPEndPoint remote) {
            try {
                rep = remote;
                socket = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);

                threadRead = new Thread(new ThreadStart(read));
                threadRead.Start();

                Connected = true;
            } catch (Exception e) {
                error(e);
            }
        }

        public void Send(byte[] data) {
            U.Send(data, data.Length);
        }

        public void Update() {
            lock (dataBuffer) {
                List<Rudp.PackageBuffer> pkgs = U.Update(dataBuffer.bt, dataBuffer.len);
                for (int i = 0; i < pkgs.Count; ++i) {
                    send(pkgs[i].buffer, pkgs[i].len);
                }
            }
        }

        public byte[] Recv() {
            return U.Recv();
        }
    }
}
