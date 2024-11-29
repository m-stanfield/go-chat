import Message, { MessageData } from "./Message"; // Import your Message component
const data = [
  { id: 1, author: "abcd1", date: 0, message: "message" },
  { id: 2, author: "abcd2", date: 1, message: "message" },
  { id: 3, author: "abcd3", date: 2, message: "message" },
  { id: 4, author: "abcd1", date: 3, message: "message" },
  { id: 5, author: "abcd2", date: 4, message: "message" },
  { id: 6, author: "abcd3", date: 5, message: "message" },
  { id: 7, author: "abcd1", date: 6, message: "message" },
  { id: 8, author: "abcd2", date: 7, message: "message" },
  { id: 9, author: "abcd3", date: 8, message: "message" },
  { id: 10, author: "abcd4", date: 9, message: "message" },
  { id: 11, author: "abcd2", date: 10, message: "message" },
  { id: 12, author: "abcd3", date: 21, message: "message" },
  { id: 13, author: "abcd4", date: 32, message: "message" },
  { id: 14, author: "abcd3", date: 23, message: "message" },
  { id: 15, author: "abcd4", date: 34, message: "message" },
] as any[];

const MockMessages = data.map((d) => new MessageData(d));
export default MockMessages;
