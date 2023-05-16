import { atom, useRecoilValue } from "recoil";
import { usersInWState } from "@src/components/sideNav1/show_workspaces";
import { UserInWorkspace } from "@src/fetchAPI/workspace";

interface SearchMemo {
    User: UserInWorkspace;
    Point: number;
}

export default function SearchUsers(input: string): UserInWorkspace[] {
    const memos: SearchMemo[] = [];
    for (const userInW of useRecoilValue(usersInWState)) {
        memos.push({
            User: userInW,
            Point: 0
        })
    }
    if (input.length == 2) {
        searchInitial(input.toLowerCase(), memos)
    }
    searchContains(input.toLowerCase(), memos)
    return createOutput(memos)
}

function createOutput(memos: SearchMemo[]): UserInWorkspace[] {
    const result: UserInWorkspace[] = [];
    for (let i = 0; i < memos.length; i ++) {
        for (let j = i + 1; j < memos.length; j ++) {
            if (memos[i] < memos[j]) {
                memos[i], memos[j] = memos[j], memos[i]
            }
        }
    }
    for (const memo of memos) {
        result.push(memo.User)
    }
    return result
}

function searchInitial(input: string, memos: SearchMemo[]):void {
    if (input.length != 2) {
        return;
    }
    for (const memo of memos) {
        const words: string[] = memo.User.name.split(' ')
        if (words.length != 2) {
            continue;
        }
        if ((words[0] == input[0] && words[1] == input[1]) || (words[0] == input[1] && words[1] == input[0])) {
            memo.Point += 5;
        }
    }
}

function searchContains(input: string, memos: SearchMemo[]):void {
    for (const memo of memos) {
        for (const word of memo.User.name) {
            if (input.length > word.length) {
                continue;
            }
            if (word.includes(input)) {
                memo.Point += 10;
                break;
            }
        }
    }
}