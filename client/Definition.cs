using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Net {
    class Definition {
        public static int BUFFER_SIZE = 1024;
    }

    class Buffer {
        public byte[] bt;
        public int len;
        public Buffer() {
            bt = new byte[Definition.BUFFER_SIZE];
            len = 0;
        }
    }
}
