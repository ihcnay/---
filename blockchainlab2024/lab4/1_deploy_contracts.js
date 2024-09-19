const Vote = artifacts.require("Vote");

module.exports = function(deployer) {

  const proposalNames = ["Proposal1", "Proposal2", "Proposal3"].map(name => web3.utils.asciiToHex(name));

  deployer.deploy(Vote, proposalNames);
};