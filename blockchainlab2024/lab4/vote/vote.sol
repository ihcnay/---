// SPDX-License-Identifier: MIT
pragma solidity >= 0.5.17;

contract Vote{
    struct Voter{
        uint weight;
        bool voted;
        address delegate;
        uint vote;
    }

    struct Proposal{
        bytes32 name;
        uint voteCount;
    }

    address public chairperson;

    mapping(address => Voter) public voters;

    Proposal[] public proposals;

    constructor(bytes32[] memory proposalNames){
        chairperson = msg.sender;
        voters[chairperson].weight = 1;
        for(uint i = 0; i < proposalNames.length;i++){
            proposals.push(Proposal({
            name: proposalNames[i],
            voteCount:0
            }));
        }
    }

    function giveRightToVote(address voter)public{
        require(msg.sender == chairperson,"Not chairperson.");
        require(!voters[voter].voted,"The voter already voted.");
        require(voters[voter].weight == 0);
        voters[voter].weight = 1;
    }

    function delegate(address to)public{        //委托投票
        Voter storage sender = voters[msg.sender];
        require(!sender.voted,"Sender already voted.");
        require(sender.weight != 0,"Sender has no right to vote.");
        require(to != msg.sender,"Cant delegate to yourself.");
        while(voters[to].delegate != address(0)){
            to = voters[to].delegate;
            require(to!= msg.sender,"Cant delegate to yourself.");
        }
        sender.delegate = to;
        sender.voted = true;
        Voter storage delegate_ = voters[to];
        if(delegate_.voted){
            proposals[delegate_.vote].voteCount += sender.weight;
        }
        else{
            delegate_.weight += sender.weight;
        }
       
    }

    function vote(uint proposal) public{
        Voter storage sender = voters[msg.sender];
        require(!sender.voted,"Already voted.");
        sender.voted = true;
        sender.vote = proposal;

        proposals[proposal].voteCount += sender.weight;
    }

    function winningProposal() public view
                returns (uint winningProposal_)
    {
        uint winningcount = 0;
        for(uint p = 0; p < proposals.length;p++){
            if(winningcount < proposals[p].voteCount){
                winningcount = proposals[p].voteCount;
                winningProposal_ = p;
            }   
        }
    }

}